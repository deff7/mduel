import Phaser from 'phaser'
import MainScene from './MainScene'

let config = {
  width: 800,
  height: 600,
  type: Phaser.AUTO,
  parent: 'canvas',
  scene: [
    MainScene,
  ],
}

let game = new Phaser.Game(config)
