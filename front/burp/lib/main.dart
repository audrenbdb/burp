import 'package:burp/home.dart';
import 'package:burp/http.dart';
import 'package:flutter/foundation.dart';
import 'package:flutter/material.dart';

void main() {
  runApp(const Burp());
}

const apiURL = kDebugMode ? "http://localhost:8080/v1" : "/v1";

class Burp extends StatelessWidget {
  const Burp({Key? key}) : super(key: key);

  @override
  Widget build(BuildContext context) {
    return MaterialApp(
      title: 'Burp',
      theme: ThemeData(
        primarySwatch: Colors.blue,
      ),
      home: Scaffold(
        body: Builder(builder: (context) {
          final api = ApiHTTP(context: context, url: apiURL);
          return Home(api: api);
        }),
      ),
    );
  }
}
