# viberBotTest

```bash
ssh -l root 95.217.131.177
tmux ls # если что-то есть, то :
tmux a

# если текущих сессий нет:
tmux 

git clone https://github.com/sergTch/viberBotTest.git
git pull

go run main.go -s release
```

Кнопки region city нельзя делать больше, чем 3 на 1 размером, так как может 
быть очень много регионов/городов и они не влезут в клавиатуру вайбера
и он не отправит сообщение
Максимальный размер клавиатуры, если не ошибаюсь, 6 колонок на 25 строк