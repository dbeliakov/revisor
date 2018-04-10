<template>
    <div id="login">
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
                <a href="#">Регистрация</a>
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
      this.$auth.login({
        data: {login: this.login, password: this.password},
        fetchUser: false
      }).then(() => {
        console.log('success')
        console.log(this.$auth)
      }).catch((err) => {
        console.log(err)
        if (!err.response.status) {
          this.error = 'Ошибка сети'
        } else if (err.response.status === 401) {
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

<style>
    body > .grid {
      height: 100%;
    }
    .column {
      max-width: 450px;
      margin: auto;
      margin-top: 40px;
    }
</style>
