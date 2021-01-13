# WebRTC Signaling Server Ayame

[![GitHub tag (latest SemVer)](https://img.shields.io/github/tag/OpenAyame/ayame.svg)](https://github.com/OpenAyame/ayame)
[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](https://opensource.org/licenses/Apache-2.0)
[![Actions Status](https://github.com/OpenAyame/ayame/workflows/Go%20Build%20&%20Format/badge.svg)](https://github.com/OpenAyame/ayame/actions)

## 時雨堂のオープンソースソフトウェアについて

利用前に https://github.com/shiguredo/oss をお読みください。

## WebRTC Signaling Server Ayame について

WebRTC Signaling Server Ayame は WebRTC 向けのシグナリングサーバです。

WebRTC の P2P でのみ動作します。また動作を 1 ルームを最大 2 名に制限することでコードを小さく保っています。

## OpenAyame プロジェクトについて

OpenAyame は WebRTC Signaling Server Ayame をオープンソースとして公開し、
継続的に開発を行うことで WebRTC をより身近に、使いやすくするプロジェクトです。

詳細については下記をご確認ください。

[OpenAyame プロジェクト](http://bit.ly/OpenAyame)

## 方針

- シグナリングの仕様の破壊的変更を可能な限り行わない
- Go のバージョンは定期的にアップデートを行う
- 依存ライブラリは定期的にアップデートを行う

## 注意

- Ayame は P2P にしか対応していません
- Ayame は 1 ルーム最大 2 名までしか対応していません

## 使ってみる

Ayame を使ってみたい人は [USE.md](docs/USE.md) をお読みください。

## Web SDK を使ってみる

[Ayame Web SDK](https://github.com/OpenAyame/ayame-web-sdk)

## Web SDK サンプルを使ってみる

[Ayame Web SDK サンプル](https://github.com/OpenAyame/ayame-web-sdk-samples)

## React サンプルを使ってみる

[OpenAyame/ayame\-react\-sample](https://github.com/OpenAyame/ayame-react-sample)

## React Native サンプルを使ってみる

[React Native WebRTC Kit のサンプルアプリケーション集](https://github.com/react-native-webrtc-kit/react-native-webrtc-kit-samples)

こちらのリポジトリの `./HelloAyame/` ディレクトリ下に Ayame の React Native サンプルがあります。

[React Native 用 WebRTC ライブラリ](https://github.com/shiguredo/react-native-webrtc-kit) を利用しています。

## 仕組みの詳細を知りたい

Ayame の詳細を知りたい人は [SPEC.md](docs/SPEC.md) をお読みください。

## Ayame Labo を使ってみる

Ayame と完全互換な STUN/TURN サーバやルーム認証を組み込んだ無料で利用可能なシグナリングサービスを時雨堂が提供しています。

[Ayame Labo](https://ayame-labo.shiguredo.jp/)

## 関連プロダクト

[hakobera/serverless-webrtc-signaling-server](https://github.com/hakobera/serverless-webrtc-signaling-server)が Ayame の互換サーバとして公開/開発されています。AWS によってサーバレスを実現した WebRTC P2P Signaling Server です。


## ライセンス

Apache License 2.0

```
Copyright 2019-2021, Shiguredo Inc.

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

