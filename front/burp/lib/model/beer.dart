class Beer {
  final String name;
  final int rating;

  Beer({required this.name, required this.rating});

  static Beer fromJson(Map<String, dynamic> payload) {
    return Beer(
      name: payload['name'] as String,
      rating: payload['rating'] as int,
    );
  }

  Map<String, dynamic> toJson() => <String, dynamic>{
      "name": name,
      "rating": rating,
  };


  Beer copyWith({String? name, int? rating}) {
    return Beer(name: name ?? this.name, rating: rating ?? this.rating);
  }
}
