import Phaser from 'phaser'

const phaserGame = new Phaser.Game(
  800,
  600,
  Phaser.AUTO,
  '',
  {
    preload: preload,
    create: create,
    update: update
  }
)

const preload = () => {}
const create = () => {}
const update = () => {}
