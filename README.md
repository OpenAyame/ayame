# WebRTC Signaling Server Ayame

[![GitHub tag (latest SemVer)](https://img.shields.io/github/tag/OpenAyame/ayame.svg)](https://github.com/OpenAyame/ayame)
[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](https://opensource.org/licenses/Apache-2.0)

## WebRTC Signaling Server Ayame について

WebRTC Signaling Server Ayame は WebRTC 向けのシグナリングサーバです。

WebRTC の P2P でのみ動作します。また動作を 1 ルームを最大 2 名に制限することでコードを小さく保っています。

AppRTC 互換のルーム機能を持っており、ルーム数はサーバスペックに依存しますが 1 万までは処理できるようにできてます。

## OpenAyame プロジェクトについて

OpenAyame プロジェクトは WebRTC Signaling Server Ayame をオープンソースとして公開し、継続的に開発を行うことで、 WebRTC を学びやすくするプロジェクトです。

詳細については下記をご確認ください。

[OpenAyame プロジェクト](http://bit.ly/OpenAyame)

## 注意

- Ayame は P2P にしか対応していません
- Ayame は 1 ルーム最大 2 名までしか対応していません
- サンプルが利用している STUN サーバは Google のものを利用しています

## 使ってみる

Ayame を使ってみたい人は [USE.md](doc/USE.md) をお読みください。

## サンプルを使ってみたい

**このリポジトリにあるサンプルと全く同じ仕組みになっています**

- Vue サンプル
    - [OpenAyame/ayame\-vue\-sample](https://github.com/OpenAyame/ayame-vue-sample)
- React サンプル
    - [OpenAyame/ayame\-react\-sample](https://github.com/OpenAyame/ayame-react-sample)

## 仕組みの詳細を知りたい

Ayame の詳細を知りたい人は [DETAIL.md](doc/DETAIL.md) をお読みください。

## Node.js (TypeScript) バージョン

**今後のメンテナンスはありません**

[OpenAyame/ayame\-nodejs](https://github.com/OpenAyame/ayame-nodejs)


## 関連プロダクト

[hakobera/serverless-webrtc-signaling-server](https://github.com/hakobera/serverless-webrtc-signaling-server)が Ayame の互換サーバとして公開/開発されています。AWS によってサーバレスを実現した WebRTC P2P Signaling Server です。


## ライセンス

Apache License 2.0

```
Copyright 2019, Shiguredo Inc, kdxu

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

WebRTC Signaling Server Ayame に関するバグ報告は GitHub Issues へお願いします。それ以外については Discord へお願いします。

### バグ報告

https://github.com/OpenAyame/ayame/issues

### Discord

ベストエフォートで運用しています。

https://discord.gg/mDesh2E

### 有料サポートについて

**時雨堂では有料サポートは提供しておりません**

- [kdxu \(Kyoko KADOWAKI\)](https://github.com/kdxu) が有料でのサポートやカスタマイズを提供しています。 Discord 経由で @kdxu へ連絡をお願いします。

