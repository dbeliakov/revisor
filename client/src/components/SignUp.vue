<template>
    <div id="login">
        <div class="column">
            <form class="ui large form" v-on:submit.prevent="signUp()">
                <div class="ui">
                <div class="field">
                    <div class="ui input">
                    <input name="first_name" v-model="first_name" placeholder="Имя" type="text">
                    </div>
                </div>
                <div class="field">
                    <div class="ui input">
                    <input name="last_name" v-model="last_name" placeholder="Фамилия" type="text">
                    </div>
                </div>
                <div class="field">
                    <div class="ui input">
                    <input name="username" v-model="username" placeholder="Логин" type="text">
                    </div>
                </div>
                <div class="field">
                    <div class="ui input">
                    <input name="password" v-model="password" placeholder="Пароль" type="password">
                    </div>
                </div>
                <button class="ui fluid large blue submit button" type="submit">Зарегистрироваться</button>
                </div>

                <div class="ui negative message" v-if="error.length > 0">{{ error }}</div>

            </form>
            <div class="ui message">
                <router-link :to="{name: 'Login'}">Вход</router-link>
            </div>
        </div>
    </div>
</template>

<script>
require('semantic-ui-css/semantic.min.css')

export default {
  name: 'SignUp',
  data () {
    return {
      first_name: '',
      last_name: '',
      username: '',
      password: '',
      error: ''
    }
  },
  methods: {
    signUp () {
      if (this.first_name.length === 0) {
        this.error = 'Имя обязательно'
        return
      }
      if (this.last_name.length === 0) {
        this.error = 'Фамилия обязательна'
        return
      }
      if (this.username.length === 0) {
        this.error = 'Логин обязателен'
        return
      }
      if (this.password.length < 6) {
        this.error = 'Пароль не короче 6 символов'
        return
      }

      this.$http.post('/auth/signup', {
        first_name: this.first_name,
        last_name: this.last_name,
        username: this.username,
        password: this.password
      }).then(() => {
        this.$router.push({name: 'Home'})
      }).catch((err) => {
        if (!err.response.status) {
          this.error = 'Ошибка сети'
        } else if (err.response.status === 406) {
          this.error = 'Неверные данные'
        } else if (err.response.status === 409) {
          this.error = 'Логин уже занят'
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
