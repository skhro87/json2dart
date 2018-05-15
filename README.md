# json2dart

Converts Json to Dart code.

_work in progress!_

### Example Usage

```
λ curl https://jsonplaceholder.typicode.com/posts/1 | go run json2dart.go
  % Total    % Received % Xferd  Average Speed   Time    Time     Time  Current
                                 Dload  Upload   Total   Spent    Left  Speed
100   292  100   292    0     0   1857      0 --:--:-- --:--:-- --:--:--  2979
2018/05/15 17:39:32 converting...
2018/05/15 17:39:32 done!
2018/05/15 17:39:32
                class Xxx {

                        final String title;
                        final String body;
                        final double userId;
                        final double id;

                        Xxx({this.title,this.body,this.userId,this.id})

                        factory Xxx.fromJson(Map<String, dynamic> json) {
                                return new Xxx(

                                        title: json['title'],
                                        body: json['body'],
                                        userId: json['userId'],
                                        id: json['id'],
                                );
                        }
                }
```

### To-Do
- [x] basic implementation
- [x] add nested object support
- [ ] add array support
- [ ] add file output support
- [ ] make it pretty ^^