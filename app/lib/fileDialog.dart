import 'package:flutter/material.dart';
import 'package:http/http.dart' as http;
import "dart:convert";

enum DialogMode {
    open,
    save
}

class DirEntry {
    final String pathUrl;
    final String name;
    final String ftype;

    const DirEntry({required this.pathUrl, required this.name, required this.ftype});

    bool empty() {
        return this.pathUrl == "" && this.name == "" && this.ftype == "";
    }

    String toString() {
        return "urlPath: ${this.pathUrl};name: ${this.name};ftype: ${this.ftype}";
    }

    bool isDir() {
        return this.ftype == "d" || this.ftype == "dir";
    }

}

class DirList {
    final String parent;
    final List<DirEntry> list;

    const DirList({required this.parent, required this.list});

    factory DirList.fromJson(Map<String, dynamic> json) {

        List<DirEntry> dirEntryList = [];

        for(final element in json["List"]) {
            DirEntry entry = DirEntry(
                pathUrl: element["PathUrl"],
                name: element["Name"],
                ftype: element["Ftype"]
            );
            dirEntryList.add(entry);
        }

        return DirList(parent: json["Parent"], list: dirEntryList);
    }
}

class DirListWidget extends StatelessWidget {
    DirList? dirList;
    void Function(BuildContext, DirEntry) pickElement;

    DirListWidget(this.pickElement, {required DirList dirList, Key? key}) : super(key : key) {
        this.dirList = dirList;
    }

    @override
    Widget build(BuildContext context) {
        return Expanded(
            child: ListView(
                children: List<Widget>.generate(
                    this.dirList!.list.length,
                    (int i) {
                        DirEntry entry = this.dirList!.list[i];
                        late Icon icon;
                        if (entry.isDir()) {
                            icon = Icon(
                                Icons.folder,
                                color: Colors.blue,
                                size: 16.0
                            );
                        } else {
                            icon = Icon(
                                Icons.text_snippet,
                                color: Colors.blue,
                                size: 16.0
                            );
                        }
                        return TextButton.icon(
                            onPressed: () {
                                pickElement(context, this.dirList!.list[i]);
                            },
                           label: Text(this.dirList!.list[i].name),
                           icon: icon
                        );
                    }
                )
            )
        );
    }
}

Future<DirList> requestDirList(String path) async {
    final Uri url = Uri(scheme: "http", host: "localhost", port: 8080, path: "/path/$path");
    final response = await http.get(url);
    debugPrint("Response: $response");
    if (response.statusCode == 200) {
        return DirList.fromJson(jsonDecode(response.body));
    } else {
        throw Exception("Failed to load: $path");
    }
}

class _FileDialog extends StatefulWidget {
    DialogMode mode;

    _FileDialog(this.mode, {Key? key}) : super(key: key);

    @override
    State<_FileDialog> createState() => _FileDialogState();
}

class _FileDialogState extends State<_FileDialog> {
    String path = "home";
    late Future<DirList> dirList;
    late DirList? currentDirList;

    @override
    initState() {
        super.initState();
        dirList = requestDirList(path);
        debugPrint("init state _FileDialogState");
    }

    void askDirList(BuildContext context, DirEntry entry) {
        if (entry.isDir()) {
            setState(() {
                path = entry.pathUrl;
                dirList = requestDirList(entry.pathUrl);
            });
        } else {
            Navigator.pop(context, entry.pathUrl);
        }
    }

    void askParent() {
        if (currentDirList == null) {
            return;
        }
        setState((){
            path = currentDirList!.parent;
            dirList = requestDirList(currentDirList!.parent);
        });
    }

    @override
    Widget build(BuildContext context) {
        debugPrint("Build fileDialogState");

        List<Widget> children = <Widget>[
            Center(child: Text(Uri.decodeFull(path))),
            TextButton.icon(
                onPressed: (){
                    if (currentDirList == null) {
                        return;
                    }
                    askParent();
                },
                icon: Icon(Icons.arrow_upward),
                label: const Text("Parent")
            ),
            FutureBuilder<DirList>(
                future: dirList,
                builder: (context, snapshot) {
                    if (snapshot.hasData) {
                        currentDirList = snapshot.data!;
                        return DirListWidget(this.askDirList, dirList: snapshot.data!);
                    } else if (snapshot.hasError) {
                        return Text("${snapshot.error}");
                    } else {
                        return const CircularProgressIndicator();
                    }
                }
            ),
        ];

        late List<Widget> finalChildren;

        if (widget.mode == DialogMode.open) {
            finalChildren = children + <Widget>[
                TextButton(
                    onPressed: (){
                        Navigator.pop(context);
                    },
                    child: Text("Close")
                ),
            ];
        } else {
            finalChildren = children + <Widget>[
                TextButton(
                    onPressed: (){
                        debugPrint("Bring route for file to save.");
                    },
                    child: Text("Save")
                ),
            ];

        }

        return Column(
            children: finalChildren
        );
    }

}

Future<String?> fileDialog(BuildContext context, DialogMode mode) {
    return showDialog<String>(
        context: context,
        builder: (BuildContext context) {
            return Dialog(child: _FileDialog(mode));
        }
    );
}
