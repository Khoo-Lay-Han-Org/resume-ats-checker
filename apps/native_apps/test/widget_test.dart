import 'package:flutter_test/flutter_test.dart';

import 'package:native_apps/main.dart';

void main() {
  testWidgets('App loads', (WidgetTester tester) async {
    await tester.pumpWidget(const App());
    expect(find.text('Resuming'), findsOneWidget);
  });
}
