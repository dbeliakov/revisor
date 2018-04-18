<template>
    <div id="login">
        <Header></Header>

        <div class="column">
            <form class="ui large form" v-on:submit.prevent="login()">
                <div class="ui">
                <div class="field">
                    <div class="ui left icon input">
                    <i class="user icon"></i>
                    <input name="username" v-model="username" placeholder="Логин" type="text">
                    </div>
                </div>
                <div class="field">
                    <div class="ui left icon input">
                    <i class="lock icon"></i>
                    <input name="password" v-model="password" placeholder="Пароль" type="password">
                    </div>
                </div>
                <button class="ui fluid large blue submit button" type="submit">Войти</button>
                </div>

                <div class="ui negative message" v-if="error.length > 0">{{ error }}</div>

            </form>
            <div class="ui message">
                <router-link :to="{name: 'SignUp'}">Регистрация</router-link>
            </div>
        </div>
    </div>
</template>

<script>
require('semantic-ui-css/semantic.min.css')

export default {
  name: 'Login',
  data () {
    return {
      username: '',
      password: '',
      error: ''
    }
  },
  methods: {
    login () {
      if (this.username.length === 0) {
        this.error = 'Логин обязателен'
        return
      }
      if (this.password.length === 0) {
        this.error = 'Пароль обязателен'
        return
      }

      this.$auth.login({
        data: {username: this.username, password: this.password}
      }).then(() => {
        console.log('success')
        console.log(this.$auth)
      }).catch((err) => {
        console.log(err)
        if (!err.response.status) {
          this.error = 'Ошибка сети'
        } else if (err.response.status === 406) {
          this.error = 'Неверный логин и/или пароль'
        } else if (err.response.status === 400) {
          this.error = 'Некорректный логин и/или пароль'
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
