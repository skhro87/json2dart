# json2dart

Converts Json to Dart code.

_work in progress!_

### Example Usage

```
λ curl https://jsonplaceholder.typicode.com/posts/1 | go run json2dart.go
  % Total    % Received % Xferd  Average Speed   Time    Time     Time  Current
                                 Dload  Upload   Total   Spent    Left  Speed
100   292  100   292    0     0   4825      0 --:--:-- --:--:-- --:--:--  4866
2018/05/15 23:43:39 converting...
2018/05/15 23:43:39 done!
2018/05/15 23:43:39
class Root {
	final String title;
	final String body;
	final double userId;
	final double id;

	Root({this.userId,this.id,this.title,this.body})

	Root.fromJson(Map<String, dynamic> json) {
		return new Root(
			userId: json['userId'],
			id: json['id'],
			title: json['title'],
			body: json['body'],
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