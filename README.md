## Необходимо написать шину.
Принимает сообщения разного типа
Проксирует их в клиентов, которые подписались на сообщения этого типа 
Хранит N последних сообщений 
Отправляет все сообщения данного типа из очереди при новом подключении клиента

### feeder
connect {stream => "in", name => $STR} # =>
msg {type => $TYPE, msg => $STRUCT} # =>

### reader
connect {stream => "out", type => $TYPE} # =>
msg {feeder => $NAME, type => $TYPE, msg => $STRUCT} # <=