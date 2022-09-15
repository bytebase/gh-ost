#!/bin/bash
set -x
table_size=100000
primary_port=
replica_port=
tps=
rate=

# first run benchmark to find out max transactions per second.
# use 50% of the max tps to simulate workload.
calc_rate() {

    rate=$(echo "$tps 0.5" | awk '{printf "%.0f", $1*$2}')
}

deploy() {
    echo "Creating replication sandbox"
    dbdeployer deploy replication 8.0.27 --nodes 2 --sandbox-directory ghostest

    echo '#!/bin/bash' >/usr/local/bin/ghostest-primary
    echo '/root/sandboxes/ghostest/m "$@"' >>/usr/local/bin/ghostest-primary
    chmod +x /usr/local/bin/ghostest-primary

    echo '#!/bin/bash' >/usr/local/bin/ghostest-replica
    echo '/root/sandboxes/ghostest/s1 "$@"' >/usr/local/bin/ghostest-replica
    chmod +x /usr/local/bin/ghostest-replica

    echo '#!/bin/bash' >/usr/local/bin/ghostest-cleanup
    echo '/root/sandboxes/ghostest/wipe_and_restart_all "$@"' >/usr/local/bin/ghostest-cleanup
    chmod +x /usr/local/bin/ghostest-cleanup
}

test_once_oltp_insert() {
    ghostest-cleanup
    ghostest-primary -uroot -e"CREATE USER IF NOT EXISTS ghost IDENTIFIED BY 'ghost'; GRANT ALL PRIVILEGES ON *.* TO ghost;"
    ghostest-primary -uroot -e"DROP DATABASE IF EXISTS db; CREATE DATABASE db;"

    primary_port=$(ghostest-primary -uroot -e "select @@port" -ss)
    replica_port=$(ghostest-replica -uroot -e "select @@port" -ss)

    sysbench \
        --mysql-host='127.0.0.1' \
        --mysql-port="$primary_port" \
        --mysql-user=ghost \
        --mysql-password=ghost \
        --mysql-db=db \
        --table-size=$table_size oltp_insert prepare

    tps=$(
        sysbench \
            --mysql-host='127.0.0.1' \
            --mysql-port="$primary_port" \
            --mysql-user=ghost \
            --mysql-password=ghost \
            --mysql-db=db \
            --table-size=$table_size --time=10 \
            oltp_insert run | grep -o 'transactions:.*)' | cut -d '(' -f 2 | cut -d ' ' -f 1
    )

    calc_rate

    sysbench \
        --mysql-host='127.0.0.1' \
        --mysql-port="$primary_port" \
        --mysql-user=ghost \
        --mysql-password=ghost \
        --mysql-db=db \
        --table-size=$table_size \
        --time=10000 \
        --rate=$rate \
        oltp_insert run &

    gh-ost \
        --execute \
        --max-load=Threads_running=25 \
        --critical-load=Threads_running=1000 \
        --assume-rbr \
        --chunk-size=10 \
        --max-lag-millis=15000 \
        --user='ghost' \
        --password='ghost' \
        --host='127.0.0.1' \
        --port="$replica_port" \
        --assume-master-host=127.0.0.1:${primary_port} \
        --database='db' \
        --table='sbtest1' \
        --verbose \
        --debug \
        --test-on-replica \
        --alter='ENGINE=InnoDB' \
        --exact-rowcount \
        --concurrent-rowcount \
        --default-retries=3 \
        --panic-flag-file=/tmp/gh-ost.panic \
        --initially-drop-old-table \
        --initially-drop-ghost-table \
        --initially-drop-socket-file \
        --serve-socket-file=/tmp/gh-ost.sock

    ghostest-replica -e"select * from db.sbtest1" -ss >/tmp/ori
    ghostest-replica -e"select * from db.\`~sbtest1_gho\`" -ss >/tmp/gho

    ori_sum="$(md5sum /tmp/ori | cut -d " " -f1)"
    gho_sum="$(md5sum /tmp/gho | cut -d " " -f1)"

    if [ "$ori_sum" != "$gho_sum" ]; then
        exit 1
    fi
}

test_once_tpcc() {
    ghostest-cleanup
    ghostest-primary -uroot -e"CREATE USER IF NOT EXISTS ghost IDENTIFIED BY 'ghost'; GRANT ALL PRIVILEGES ON *.* TO ghost;"
    ghostest-primary -uroot -e"DROP DATABASE IF EXISTS db; CREATE DATABASE db;"

    primary_port=$(ghostest-primary -uroot -e "select @@port" -ss)
    replica_port=$(ghostest-replica -uroot -e "select @@port" -ss)

    /tpcc.lua prepare --mysql-host='127.0.0.1' \
        --mysql-port="$primary_port" \
        --mysql-user=ghost \
        --mysql-password=ghost \
        --mysql-db=db \
        --use-fk=0 \
        --scale=1

    tps=$(
        /tpcc.lua run \
            --mysql-host='127.0.0.1' \
            --mysql-port="$primary_port" \
            --mysql-user=ghost \
            --mysql-password=ghost \
            --mysql-db=db \
            --time=10 \
            --scale=1 \
            | grep -o 'transactions:.*)' | cut -d '(' -f 2 | cut -d ' ' -f 1
    )

    calc_rate

    /tpcc.lua run \
        --mysql-host='127.0.0.1' \
        --mysql-port="$primary_port" \
        --mysql-user=ghost \
        --mysql-password=ghost \
        --mysql-db=db \
        --time=10000 \
        --rate=$rate \
        &

    gh-ost \
        --execute \
        --max-load=Threads_running=25 \
        --critical-load=Threads_running=1000 \
        --assume-rbr \
        --chunk-size=10 \
        --max-lag-millis=15000 \
        --user='ghost' \
        --password='ghost' \
        --host='127.0.0.1' \
        --port="$replica_port" \
        --assume-master-host=127.0.0.1:${primary_port} \
        --database='db' \
        --table='order_line1' \
        --verbose \
        --debug \
        --test-on-replica \
        --alter='ENGINE=InnoDB' \
        --exact-rowcount \
        --concurrent-rowcount \
        --default-retries=3 \
        --panic-flag-file=/tmp/gh-ost.panic \
        --initially-drop-old-table \
        --initially-drop-ghost-table \
        --initially-drop-socket-file \
        --serve-socket-file=/tmp/gh-ost.sock

    ghostest-replica -e"select * from db.order_line1" -ss >/tmp/ori
    ghostest-replica -e"select * from db.\`~order_line1_gho\`" -ss >/tmp/gho

    ori_sum="$(md5sum /tmp/ori | cut -d " " -f1)"
    gho_sum="$(md5sum /tmp/gho | cut -d " " -f1)"

    if [ "$ori_sum" != "$gho_sum" ]; then
        exit 1
    fi
}

deploy
test_once_oltp_insert
test_once_tpcc
exit 0
