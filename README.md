# json2dart

Converts Json to Dart code.

_work in progress!_

### Example Usage

```javascript
λ cat example.json
{
  "id": 55,
  "name": "Alf",
  "something": {
    "beard": "long",
    "items": 1
  },
  "items": [
    {
      "score": 5
    }
  ]
}
```

```dart
λ go build && cat example.json | json2dart.exe
2018/05/16 11:33:37 converting...
2018/05/16 11:33:37 done!
2018/05/16 11:33:37
class Something {
        final String beard;
        final num items;

        Something({this.beard,this.items})

        Something.fromJson(Map<String, dynamic> json) {
                return new Something(
                        beard: json['beard'],
                        items: json['items'],
                );
        }
}

class Item {
        final num score;

        Item({this.score})

        Item.fromJson(Map<String, dynamic> json) {
                return new Item(
                        score: json['score'],
                );
        }
}

class Root {
        final num id;
        final String name;
        final Something something;
        final List<Item> items;

        Root({this.items,this.id,this.name,this.something})

        Root.fromJson(Map<String, dynamic> json) {
                return new Root(
                        name: json['name'],
                        something: json['something'],
                        items: json['items'],
                        id: json['id'],
                );
        }
}
```

### To-Do
- [x] basic implementation
- [x] add nested object support
- [x] add array support
- [ ] improve input/output options
- [ ] proper tests / edge cases
- [ ] make it pretty ^^
