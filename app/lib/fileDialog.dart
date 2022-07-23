import 'package:flutter/material.dart';
import 'package:http/http.dart' as http;

class _FileDialog extends StatefulWidget {
    _FileDialog({Key? key}) : super(key: key);

    @override
    State<_FileDialog> createState() => _FileDialogState();
}

class _FileDialogState extends State<_FileDialog> {

    @override
    Widget build(BuildContext context) {

        List<Widget> children = <Widget>[
            Center(child: Text("some path")),
            Expanded(
                child: ListView()
            ),
            Row(
                children: <Widget>[
                    TextButton(
                        onPressed: (){},
                        child: Text("Ok")
                    ),
                    TextButton(
                        onPressed: (){
                            Navigator.pop(context);
                        },
                        child: Text("Cancel")
                    )
                ]
            ),
        ];

        return Column(
            children: children
        );
    }
}

Future<void> fileDialog(BuildContext context, String path) async {

    await showDialog<void>(
        context: context,
        builder: (BuildContext context) {
            return Dialog(child: _FileDialog());
        }
    );

    debugPrint("Hello! $path");
}
