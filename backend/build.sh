set -e
GOOS=linux GOARCH=amd64 go build

host="root@..."

ssh $host 'systemctl stop epaper-dashboard.service'
scp epaper-dashboard $host:~/dev/epaper-dashboard
ssh $host 'systemctl restart epaper-dashboard.service'
