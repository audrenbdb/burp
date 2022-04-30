import 'dart:async';

import 'package:burp/http.dart';
import 'package:burp/model/beer.dart';
import 'package:flutter/material.dart';
import 'package:flutter_rating_bar/flutter_rating_bar.dart';

class Home extends StatefulWidget {
  final ApiHTTP api;
  const Home({Key? key, required this.api}) : super(key: key);

  @override
  State<Home> createState() => _HomeState();
}

class _HomeState extends State<Home> {
  final _beers = StreamController<List<Beer>>();
  final BeerFilter _filter = BeerFilter();
  Timer? debounce;

  fetchBeers() async {
    _beers.add(await widget.api.fetchBeers(filter: _filter));
  }

  setNameFilter(String? name) {
    debounce?.cancel();
    debounce = Timer(const Duration(milliseconds: 350), () {
      _filter.nameContains = name;
      fetchBeers();
    });
  }

  deleteBeer(Beer b) async {
    await widget.api.deleteBeer(beer: b);
    fetchBeers();
  }

  Function(double) setBeerRating(List<Beer> beers, Beer beer) {
    return (double rating) async {
      final r = rating.toInt();
      final updatedBeer = beer.copyWith(rating: r);
      await widget.api.saveBeer(beer: updatedBeer);
      fetchBeers();
    };
  }

  Function(Beer) addBeer(List<Beer> beers) => (Beer beer) async {
        await widget.api.saveBeer(beer: beer);
        fetchBeers();
      };

  @override
  void initState() {
    super.initState();
    fetchBeers();
  }

  @override
  void dispose() {
    _beers.close();
    super.dispose();
  }

  // This widget is the root of your application.
  @override
  Widget build(BuildContext context) {
    return StreamBuilder(
        stream: _beers.stream,
        builder: (BuildContext context, AsyncSnapshot<List<Beer>> snapshot) {
          if (!snapshot.hasData) {
            return const Center(child: CircularProgressIndicator());
          }
          final beers = snapshot.data!;
          return Scaffold(
            appBar: AppBar(
              // The search area here
                title: Container(
                  width: double.infinity,
                  height: 40,
                  decoration: BoxDecoration(
                      color: Colors.white, borderRadius: BorderRadius.circular(5)),
                  child: Center(
                    child: TextField(
                      onChanged: (String value) => setNameFilter(value.isEmpty ? null : value),
                      decoration: const InputDecoration(
                          prefixIcon:Icon(Icons.search),
                          hintText: 'Search a beverage...',
                          border: InputBorder.none),
                    ),
                  ),
                )),
            body: ListView.builder(
                itemCount: beers.length,
                itemBuilder: (context, index) {
                  final beer = beers[index];
                  return Dismissible(
                    key: ValueKey<Beer>(beer),
                    onDismissed: (_) => deleteBeer(beer),
                    background: Container(color: Colors.red),
                    child: ListTile(
                        title: Text(beer.name),
                        trailing: RatingBar.builder(
                            initialRating: beer.rating.toDouble(),
                            minRating: 0,
                            maxRating: 4,
                            itemCount: 4,
                            itemBuilder: (context, _) => const Icon(
                                  Icons.star,
                                  color: Colors.amber,
                                ),
                            onRatingUpdate: setBeerRating(beers, beer))),
                  );
                }),
            floatingActionButton: FloatingActionButton(
              child: const Icon(Icons.add),
              onPressed: _navigateToBeerForm(
                  context: context, onSubmit: addBeer(beers)),
            ),
          );
        });
  }

  Function() _navigateToBeerForm(
          {required BuildContext context, required Function(Beer) onSubmit}) =>
      () {
        Navigator.push(
            context,
            MaterialPageRoute(
                builder: (context) => BeerForm(
                      onSubmit: onSubmit,
                    )));
      };
}

class BeerForm extends StatefulWidget {
  final Function(Beer) onSubmit;

  const BeerForm({Key? key, required this.onSubmit}) : super(key: key);

  @override
  State<BeerForm> createState() => _BeerFormState();
}

class _BeerFormState extends State<BeerForm> {
  final _formKey = GlobalKey<FormState>();
  final name = TextEditingController();
  int rating = 0;

  @override
  void dispose() {
    // TODO: implement dispose
    super.dispose();
    name.dispose();
  }

  String? _validateName(String? name) {
    if (name == null) {
      return "beer name cannot be empty";
    }
    if (name.length < 2) {
      return "beer name must have at least 2 characters";
    }
    if (name.length > 15) {
      return "beer name must have maximum 15 characters";
    }
    return null;
  }

  _setRating(double r) => rating = r.toInt();

  _submit() async {
    await widget.onSubmit(Beer(name: name.text, rating: rating));
    Navigator.of(context).pop();
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
        appBar: AppBar(title: const Text("Add new beer")),
        body: Form(
            key: _formKey,
            child: Padding(
              padding: const EdgeInsets.all(36.0),
              child: Column(
                children: [
                  TextFormField(
                      autofocus: true,
                      controller: name,
                      decoration: const InputDecoration(
                        labelText: "Ex.: Guinness",
                      ),
                      validator: _validateName,
                  ),
                  const Padding(padding: EdgeInsets.all(24)),
                  RatingBar.builder(
                      initialRating: rating.toDouble(),
                      minRating: 0,
                      maxRating: 4,
                      itemCount: 4,
                      itemBuilder: (context, _) => const Icon(
                        Icons.star,
                        color: Colors.amber,
                      ),
                      onRatingUpdate: _setRating),
                  const Padding(padding: EdgeInsets.all(24)),
                  ElevatedButton(
                      onPressed: () {
                        if (_formKey.currentState!.validate()) {
                          _submit();
                        }
                      },
                      child: const Padding(
                        padding: EdgeInsets.only(top: 12, bottom: 12, left: 24, right: 24),
                        child: Text('SUBMIT'),
                      ),
                    ),
                ],
              ),
            )));
  }
}
