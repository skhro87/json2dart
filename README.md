# json2dart

Converts Json to Dart code.

_work in progress!_

### Help

```
λ ./json2dart --help
NAME:
   json2dart - convert json to dart code

USAGE:
   json2dart [global options] command [command options] [arguments...]

VERSION:
   0.1.0

COMMANDS:
     help, h  Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --in value     specify input file
   --out value    specify output file(s) folder (default: "./model_generated")
   --class value  name of root class (default: "Root")
   --split        split generated classes into multiple files
   --no-files     disable output to file(s)
   --help, -h     show help
   --version, -v  print the version
```

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
λ cat example.json | json2dart
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
