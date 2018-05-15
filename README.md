# json2dart

Converts Json to Dart code.

_work in progress!_

### Example Usage

```
λ curl https://jsonplaceholder.typicode.com/posts/1 | go run json2dart.go
  % Total    % Received % Xferd  Average Speed   Time    Time     Time  Current
                                 Dload  Upload   Total   Spent    Left  Speed
100   292  100   292    0     0   4336      0 --:--:-- --:--:-- --:--:--  4358
2018/05/15 23:41:49 converting...
2018/05/15 23:41:49 done!
2018/05/15 23:41:49
class Xxx {
	final String body;
	final double userId;
	final double id;
	final String title;

	Xxx({this.title,this.body,this.userId,this.id})

	Xxx.fromJson(Map<String, dynamic> json) {
		return new Xxx(
			body: json['body'],
			userId: json['userId'],
			id: json['id'],
			title: json['title'],
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