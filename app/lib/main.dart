import 'package:flutter/material.dart';
import "dart:convert";

import 'package:http/http.dart' as http;

void main() {
  runApp(const MyApp());
}

class DirEntry {
    String urlPath = "";
    String name = "";
    String ftype = "";

    DirEntry(this.urlPath, this.name, this.ftype);

    bool empty() {
        return this.urlPath == "" && this.name == "" && this.ftype == "";
    }

    String toString() {
        return "urlPath: ${this.urlPath};name: ${this.name};ftype: ${this.ftype}";
    }

    bool isDir() {
        return this.ftype == "d" || this.ftype == "dir";
    }

}

class DirList {
    String parent = "";
    List<DirEntry> list = [];
}

class MyApp extends StatelessWidget {
  const MyApp({Key? key}) : super(key: key);

  // This widget is the root of your application.
  @override
  Widget build(BuildContext context) {

    return MaterialApp(
      title: 'OTG Site Builder',
      theme: ThemeData(
        // This is the theme of your application.
        //
        // Try running your application with "flutter run". You'll see the
        // application has a blue toolbar. Then, without quitting the app, try
        // changing the primarySwatch below to Colors.green and then invoke
        // "hot reload" (press "r" in the console where you ran "flutter run",
        // or simply save your changes to "hot reload" in a Flutter IDE).
        // Notice that the counter didn't reset back to zero; the application
        // is not restarted.
        primarySwatch: Colors.blue,
      ),
      home: const WebsiteScreen(),
    );
  }
}

Widget buildTextEntry(String label) {
    return Container(
        child: TextFormField(
            decoration: InputDecoration(
                labelText: label
            )
        )
    , height: 64, width: 240);
}

class WebsiteData extends StatefulWidget {
    const WebsiteData({Key? key}) : super(key : key);

    @override
    State<WebsiteData> createState() => _WebsiteDataState();
}

class _WebsiteDataState extends State<WebsiteData> {
    final GlobalKey<FormState> _formKey = GlobalKey<FormState>();

    late Container title;
    late Container output;
    late Container postsPerPage;

    @override
    Widget build(BuildContext context) {

        title = buildTextEntry("Title") as Container;
        postsPerPage = buildTextEntry("Posts per page") as Container;
        output = buildTextEntry("Output") as Container;

        return Form(
            key: _formKey,
            child: Column(children: <Widget>[
                title,
                postsPerPage,
                output,
                Row(children: [
                    TextButton(
                        child: const Text("Save"),
                        onPressed: () {debugPrint("Save");}
                    ),
                    TextButton(
                        child: const Text("Reset"),
                        onPressed: () {
                            _formKey.currentState!.reset();
                            debugPrint("Reset");
                        }
                    )
                ])
            ])
        );

    }

}

Future<void> loadWebsite(BuildContext context, {String path = "home"}) async {
    final String finalPath = "path/$path";

    var url = Uri(scheme: "http", host: "localhost", port: 8080, path: finalPath);
    debugPrint("Uri: $url");
    DirEntry? entry = DirEntry("", "", "");

    try {
        while(entry != null && (entry.empty() || entry.isDir())) {

        late http.Response result;

        if (entry.isDir()) {
            final url = Uri(scheme: "http", host: "localhost", port: 8080, path: "/path/${entry.urlPath}");
            result = await http.get(url);
        } else {
            result = await http.get(url);
        }
        final data = jsonDecode(result.body);
        debugPrint("${result.statusCode}");
        debugPrint("${data.runtimeType}");

        String parent = data["Parent"];

        List<Widget> widgetDirEntry = [
            SimpleDialogOption(
                onPressed: () {
                    DirEntry entry = DirEntry(parent, ".. (Parent)", 'd');
                    Navigator.pop(context, entry);
                },
                child: const Text(".. (Parent)")
            )
        ];

        for (final child in data["List"]){
            // debugPrint("child: $child; ${child['Name']}");
            SimpleDialogOption option = SimpleDialogOption(
                onPressed: () {
                    DirEntry entry = DirEntry(child["PathUrl"], child["Name"], child["Ftype"]);
                    Navigator.pop(context, entry);
                },
                child: Text(child["Name"])
            );

            widgetDirEntry.add(option);
        }

        final String pathText = entry.isDir() ? entry.name : path;

        entry = await showDialog<DirEntry>(
            context: context,
            builder: (BuildContext context) {
                return SimpleDialog(
                    title: Text(pathText),
                    children: widgetDirEntry
                );
            }
        );

        debugPrint("entry: $entry");

        } //while

    } catch (err) {
        debugPrint("error: $err");
    }
}

class WebsiteScreen extends StatelessWidget {
    const WebsiteScreen({Key? key}) : super(key : key);

    @override
    Widget build(BuildContext context) {
        return Scaffold(
            appBar: AppBar(
                title: const Text('OTG Site Builder')
            ),
            body: Row(children: <Widget>[
                Column(children:<Widget>[
                    TextButton(
                        onPressed: () {
                            debugPrint("New website");
                        },
                        child: const Text("New website")),
                    TextButton(
                        onPressed: () {
                            debugPrint("Load website");
                            loadWebsite(context);
                            debugPrint("");
                        },
                        child: const Text("Load website")),
                    TextButton(
                        onPressed: () {
                            debugPrint("Build website");
                        },
                        child: const Text("Build website")),
                ]),
                WebsiteData()
            ])
        );
    }

}

class MyHomePage extends StatefulWidget {
  const MyHomePage({Key? key, required this.title}) : super(key: key);

  // This widget is the home page of your application. It is stateful, meaning
  // that it has a State object (defined below) that contains fields that affect
  // how it looks.

  // This class is the configuration for the state. It holds the values (in this
  // case the title) provided by the parent (in this case the App widget) and
  // used by the build method of the State. Fields in a Widget subclass are
  // always marked "final".

  final String title;

  @override
  State<MyHomePage> createState() => _MyHomePageState();
}

class _MyHomePageState extends State<MyHomePage> {
  int _counter = 0;

  void _incrementCounter() {
    setState(() {
      // This call to setState tells the Flutter framework that something has
      // changed in this State, which causes it to rerun the build method below
      // so that the display can reflect the updated values. If we changed
      // _counter without calling setState(), then the build method would not be
      // called again, and so nothing would appear to happen.
      _counter++;
    });
  }

  @override
  Widget build(BuildContext context) {
    // This method is rerun every time setState is called, for instance as done
    // by the _incrementCounter method above.
    //
    // The Flutter framework has been optimized to make rerunning build methods
    // fast, so that you can just rebuild anything that needs updating rather
    // than having to individually change instances of widgets.
    return Scaffold(
      appBar: AppBar(
        // Here we take the value from the MyHomePage object that was created by
        // the App.build method, and use it to set our appbar title.
        title: Text(widget.title),
      ),
      body: Center(
        // Center is a layout widget. It takes a single child and positions it
        // in the middle of the parent.
        child: Column(
          // Column is also a layout widget. It takes a list of children and
          // arranges them vertically. By default, it sizes itself to fit its
          // children horizontally, and tries to be as tall as its parent.
          //
          // Invoke "debug painting" (press "p" in the console, choose the
          // "Toggle Debug Paint" action from the Flutter Inspector in Android
          // Studio, or the "Toggle Debug Paint" command in Visual Studio Code)
          // to see the wireframe for each widget.
          //
          // Column has various properties to control how it sizes itself and
          // how it positions its children. Here we use mainAxisAlignment to
          // center the children vertically; the main axis here is the vertical
          // axis because Columns are vertical (the cross axis would be
          // horizontal).
          mainAxisAlignment: MainAxisAlignment.center,
          children: <Widget>[
            const Text(
              'You have pushed the button this many times:',
            ),
            Text(
              '$_counter',
              style: Theme.of(context).textTheme.headline4,
            ),
          ],
        ),
      ),
      floatingActionButton: FloatingActionButton(
        onPressed: _incrementCounter,
        tooltip: 'Increment',
        child: const Icon(Icons.add),
      ), // This trailing comma makes auto-formatting nicer for build methods.
    );
  }
}
