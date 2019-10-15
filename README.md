# WebRTC Signaling Server Ayame

[![GitHub tag (latest SemVer)](https://img.shields.io/github/tag/OpenAyame/ayame.svg)](https://github.com/OpenAyame/ayame)
[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](https://opensource.org/licenses/Apache-2.0)
[![Actions Status](https://github.com/OpenAyame/ayame/workflows/Go%20Build%20&%20Format/badge.svg)](https://github.com/OpenAyame/ayame/actions)

## WebRTC Signaling Server Ayame について

WebRTC Signaling Server Ayame は WebRTC 向けのシグナリングサーバです。

WebRTC の P2P でのみ動作します。また動作を 1 ルームを最大 2 名に制限することでコードを小さく保っています。

AppRTC 互換のルーム機能を持っており、ルーム数はサーバスペックに依存しますが 1 万までは処理できるようにできてます。

## OpenAyame プロジェクトについて

OpenAyame プロジェクトは WebRTC Signaling Server Ayame をオープンソースとして公開し、継続的に開発を行うことで、 WebRTC を学びやすくするプロジェクトです。

詳細については下記をご確認ください。

[OpenAyame プロジェクト](http://bit.ly/OpenAyame)

## 開発について

Ayame はオープンソースソフトウェアですが、開発についてはオープンではありません。
そのためコメントやプルリクエストを頂いてもすぐには採用はしません。

まずは Discord にてご連絡ください。

## 注意

- Ayame は P2P にしか対応していません
- Ayame は 1 ルーム最大 2 名までしか対応していません
- サンプルが利用している STUN サーバは Google のものを利用しています

## 使ってみる

Ayame を使ってみたい人は [USE.md](doc/USE.md) をお読みください。

## SDK を使ってみる

Ayame は Web SDK と Android SDK を提供しています。現在 iOS SDK を開発中です。

- [Ayame Web SDK](https://github.com/OpenAyame/ayame-web-sdk)
    - [Ayame Web SDK サンプル](https://github.com/OpenAyame/ayame-web-sdk-samples)
- [Ayame Android SDK](https://github.com/OpenAyame/ayame-android-sdk)
    - [Ayame Android SDK サンプル](https://github.com/OpenAyame/ayame-android-sdk-samples)
- [Ayame iOS SDK](https://github.com/OpenAyame/ayame-ios-sdk)
    - 開発中です

## React サンプルを使ってみる

**このリポジトリにあるサンプルと全く同じ動作になっています**

[OpenAyame/ayame\-react\-sample](https://github.com/OpenAyame/ayame-react-sample)

## React Native サンプルを使ってみる

[Ayame React Native サンプル](https://github.com/OpenAyame/ayame-react-native-sample)

[React Native 用 WebRTC ライブラリ](https://github.com/shiguredo/react-native-webrtc-kit) を利用しています。

## WebRTC シグナリングサービス Ayame Lite を使ってみる

Ayame を利用した無料で TURN サーバまで利用可能なシグナリングサービスを提供しています。

[Ayame Lite (オープンベータ)](https://ayame-lite.shiguredo.jp/beta)

## 仕組みの詳細を知りたい

Ayame の詳細を知りたい人は [DETAIL.md](doc/DETAIL.md) をお読みください。

## 関連プロダクト

[hakobera/serverless-webrtc-signaling-server](https://github.com/hakobera/serverless-webrtc-signaling-server)が Ayame の互換サーバとして公開/開発されています。AWS によってサーバレスを実現した WebRTC P2P Signaling Server です。

## ライセンス

Apache License 2.0

```
Copyright 2019, Shiguredo Inc, Kyoko Kadowaki (kdxu)

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
```

## サポートについて

### 要望や PR について

まずは Discord へお願いします。

### バグ報告

WebRTC Signaling Server Ayame に関するバグ報告は GitHub Issues へお願いします。

https://github.com/OpenAyame/ayame/issues

### Discord

ベストエフォートで運用しています。

https://discord.gg/mDesh2E

### 有料サポートについて

**時雨堂では有料サポートは提供しておりません**

- [kdxu \(Kyoko KADOWAKI\)](https://github.com/kdxu) が有料でのサポートやカスタマイズを提供しています。 Discord 経由で @kdxu へ連絡をお願いします。

