import 'phaser'

const state_fetch_id = "Получение ID"
const state_wait_players = "Ждем соперника"
const state_wait_game = "Начинаем"
const state_game_running = "Игра в процессе"

class MainScene extends Phaser.Scene {
  constructor() {
    super({
      key: 'MainScene'
    })
  }

  preload() {
    this.conn = new WebSocket('ws://localhost:3000/socket')
    this.clientState = state_fetch_id
    this.conn.onmessage = this.processMessage.bind(this)
  }

  resetInput() {
    this.currentInput = ''
  }

  processMessage(evt) {
    let data = JSON.parse(evt.data)
    if (this.clientState === state_fetch_id) {
      this.playerID = data.PlayerID
      this.playerIndex = data.PlayerIndex
      this.clientState = state_wait_players
    } else if (this.clientState === state_wait_players) {
      if (data.Ready == true) {
        this.clientState = state_wait_game
      }
    } else if (this.clientState === state_wait_game) {
      if (data.Players !== undefined) {
        this.state = data
        this.clientState = state_game_running
      }
    } else if (this.clientState === state_game_running) {
      this.state = data
    }
  }

  sendMessage(suffix) {
    if (this.keyLatch) {
      return
    }

    let input = this.currentInput
    if (suffix !== undefined) {
      input += suffix
    }

    let data = {
      PlayerID: this.playerID,
      Input: input,
    }
    this.conn.send(JSON.stringify(data))
    this.resetInput() 
    this.keyLatch = true
  }

  myState() {
    return this.state.Players[this.playerIndex]
  }

  enemyState() {
    return this.state.Players[1 - this.playerIndex]
  }

  createKeyboard() {
    this.keyboardMap = {}
    let keyA = Phaser.Input.Keyboard.KeyCodes.A,
      keyZero = Phaser.Input.Keyboard.KeyCodes.ZERO
    let codeA = 'a'.charCodeAt(0),
      codeZero = '0'.charCodeAt(0)

    for (let c = keyA;
      c <= Phaser.Input.Keyboard.KeyCodes.Z;
      c++) {
      this.keyboardMap[String.fromCharCode(codeA + c - keyA)] = c
    }

    for (let c = keyZero;
      c <= Phaser.Input.Keyboard.KeyCodes.NINE;
      c++) {
      this.keyboardMap[String.fromCharCode(codeZero + c - keyZero)] = c
    }
    this.keyboard = this.input.keyboard.addKeys(this.keyboardMap)
    this.keyPressed = {}
  }

  createControlKeys() {
    this.empowerKey = this.input.keyboard.addKey(Phaser.Input.Keyboard.KeyCodes.SPACE)
    this.castKey = this.input.keyboard.addKey(Phaser.Input.Keyboard.KeyCodes.ENTER)
    this.backspaceKey = this.input.keyboard.addKey(Phaser.Input.Keyboard.KeyCodes.BACKSPACE)
    this.keyLatch = false
  }

  create() {
    let hstep = 30
    this.texts = {
      clientState: this.add.text(10, hstep, '', { font: '12px Arial', fill: '#00ff00' }),
      gameState: this.add.text(10, hstep*2, '', { font: '12px Arial', fill: '#00ff00' }),
      playerID: this.add.text(10, hstep*3, '', { font: '12px Arial', fill: '#00ff00' }),
      input: this.add.text(10, hstep*4, '', { font: '20px Arial', fill: '#ffffff' }),
    }
    this.createKeyboard()
    this.createControlKeys()
    this.resetInput()
  }

  createBolt(byEnemy) {
    let bolt = this.add.graphics()
    if (this.bolts === undefined) {
      this.bolts = {}
    }
    if (byEnemy) {
      if (this.bolts.enemy !== undefined) {
        return
      }
      this.bolts.enemy = {}
      this.bolts.enemy.distance = this.enemyState().Spell.Distance
      this.bolts.enemy.speed = this.enemyState().Spell.BoltSpeed

      bolt.fillStyle(0xff0000)
      bolt.fillCircle(15, 15, 10)
      this.bolts.enemy.obj = bolt
    } else {
      if (this.bolts.my !== undefined) {
        return
      }
      this.bolts.my = {}
      this.bolts.my.distance = this.myState().Spell.Distance
      this.bolts.my.speed = this.myState().Spell.BoltSpeed

      bolt.fillStyle(0x00ff00)
      bolt.fillCircle(15, 500, 10)
      this.bolts.my.obj = bolt
    }
  }

  destroyBolt(enemy) {
    if (this.bolts === undefined) {
      return
    }
    if (enemy) {
      if (this.bolts.enemy !== undefined) {
        this.bolts.enemy.obj.destroy()
        delete this.bolts.enemy
      }
    } else if (this.bolts.my !== undefined) {
      this.bolts.my.obj.destroy()
      delete this.bolts.my
    }
  }

  checkCastedSpells() {
    if (this.myState().Spell.BoltSpeed > 0) {
      this.createBolt()
    } else {
      this.destroyBolt()
    }

    if (this.enemyState().Spell.BoltSpeed > 0) {
      this.createBolt(true)
    } else {
      this.destroyBolt(true)
    }
  }

  update(time, delta) {
    this.texts.clientState.setText(this.clientState)
    if (this.clientState === state_game_running) {
      this.texts.gameState.setText(JSON.stringify(this.myState()))
      for(let k in this.keyboard) {
        if (this.keyboard[k].isDown) {
          if (!(k in this.keyPressed)) {
            this.currentInput += k
          }
          this.keyPressed[k] = true
        } else {
          delete this.keyPressed[k]
        }
      }

      if (this.empowerKey.isDown) {
        this.sendMessage(' ')
      } else if (this.castKey.isDown) {
        this.sendMessage()
      } else if (this.backspaceKey.isDown) {
        if (!this.keyLatch) {
          this.currentInput = this.currentInput.slice(0, -1)
          this.keyLatch = true
        }
      } else {
        this.keyLatch = false
      }

      this.texts.input.setText(`Текущий ввод: ${this.currentInput}`)

      this.checkCastedSpells()
    }
    this.texts.playerID.setText(`Твой ID: ${this.playerID}`)

    if (this.bolts !== undefined) {
      if (this.bolts.enemy !== undefined) {
        this.bolts.enemy.obj.y += this.bolts.enemy.speed * 0.07
      }
      if (this.bolts.my !== undefined) {
        this.bolts.my.obj.y -= this.bolts.my.speed * 0.07
      }
    }
  }
}

export default MainScene
