import 'dart:convert';

import 'package:burp/model/beer.dart';
import 'package:flutter/material.dart';
import 'package:http/http.dart' as http;

class BeerFilter {
  String? nameContains;
  List<int>? ratings;

  BeerFilter({this.nameContains, this.ratings});

  Map<String, String>? toURLParams() {
    if (nameContains == null && ratings == null) {
      return null;
    }
    final params = <String, String>{};
    if (nameContains != null) {
      params["name"] = nameContains!;
    }
    if (ratings != null) {
      params["ratings"] = ratings!.join(",");
    }
    return params;
  }
}

class ApiHTTP {
  final BuildContext context;
  final String url;

  ApiHTTP({required this.context, required this.url});

  Future<List<Beer>> fetchBeers({BeerFilter? filter}) async {
    final params = filter?.toURLParams();
    final body = await _get(endpoint: "/beers", params: params);
    final beers = jsonDecode(body) as List;
    return beers.map((b) => Beer.fromJson(b)).toList();
  }

  Future<void> saveBeer({required Beer beer}) async {
    await _put(endpoint: "/beers/${beer.name}", resource: beer.toJson());
  }

  Future<void> deleteBeer({required Beer beer}) async {
    await _delete(endpoint: "/beers/${beer.name}");
  }

  Future<String> _get(
      {required String endpoint, Map<String, String>? params}) async {
    final uri = _buildURI(endpoint: endpoint, params: params);
    return await _handleFetch(() => http.get(uri));
  }

  Future<String> _put({required String endpoint, required Map<String, dynamic> resource}) async {
    final uri = _buildURI(endpoint: endpoint);
    return await _handleFetch(() => http.put(uri, body: jsonEncode(resource), headers: {
      "Content-Type": "application/json",
    }));
  }

  Future<String> _delete({required String endpoint}) async {
    final uri = _buildURI(endpoint: endpoint);
    return await _handleFetch(() => http.delete(uri));
  }

  Future<String> _handleFetch(Function fetch) async {
    try {
      final resp = await fetch() as http.Response;
      if (resp.statusCode ~/ 100 != 2) {
        throw resp.body;
      }
      return resp.body;
    } catch(e) {
      ScaffoldMessenger.of(context).showSnackBar(SnackBar(
        content: Text(e.toString()),
      ));
      rethrow;
    }
  }

  Uri _buildURI({required String endpoint, Map<String, String>? params}) {
    endpoint = url + endpoint;
    if (params != null) {
      String queryString = Uri(queryParameters: params).query;
      endpoint += "?" + queryString;
    }
    return Uri.parse(endpoint);
  }
}
