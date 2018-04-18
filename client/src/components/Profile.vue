<template>
  <div id="profile">
  <h1>Смена пароля</h1>
  <div class="column">
    <form class="ui large form" v-on:submit.prevent="changePassword()">
      <div class="ui">
        <div class="field">
          <div class="ui input">
            <input name="password" v-model="old_password" placeholder="Текущий пароль" type="password">
          </div>
        </div>
        <div class="field">
          <div class="ui input">
            <input name="password" v-model="new_password" placeholder="Новый пароль" type="password">
          </div>
        </div>
        <button class="ui fluid large blue submit button" type="submit">Сменить пароль</button>
      </div>
      <div class="ui negative message" v-if="error.length > 0">{{ error }}</div>
    </form>
  </div>
  </div>
</template>

<script>
export default {
  name: 'Profile',
  data () {
    return {
      old_password: '',
      new_password: '',
      error: ''
    }
  },
  methods: {
    changePassword () {
      if (this.old_password.length === 0) {
        this.error = 'Текущий пароль необходим'
        return
      }
      if (this.new_password.length < 6) {
        this.error = 'Новый пароль должен быть не короче 6 символов'
        return
      }
      this.$http.post('/auth/change/password', {
        old_password: this.old_password,
        new_password: this.new_password
      }).then(() => {
        this.$router.push({path: '/'})
      }).catch((err) => {
        if (!err.response.status) {
          this.error = 'Ошибка сети'
        } else if (err.response.status === 406) {
          this.error = 'Некорректный пароль'
        } else {
          this.error = 'Неизвестная ошибка'
        }
      })
    }
  }
}
</script>

<style scoped>
  body > .grid {
    height: 100%;
  }
  .column {
    max-width: 450px;
    margin: auto;
    margin-top: 40px;
  }
</style>
