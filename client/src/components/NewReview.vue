<template>
  <div id="new_review">
    <h1>Новое ревью</h1>
    <div class="column">
      <form class="ui large form" v-on:submit.prevent="createNewReview()">
        <div class="ui">
          <div class="field">
            <div class="ui input">
              <input name="name" v-model="name" placeholder="Заголовок" type="text">
            </div>
          </div>
          <div class="field">
            <div class="ui input">
              <input name="reviewers" v-model="reviewers" placeholder="Ревьюеры" type="text">
            </div>
          </div>
          <div class="field">
            <label for="file" class="ui icon button">
              <i class="file icon"></i>
              {{ filename }}</label>
            <input type="file" id="file" style="display:none" v-on:change="changeFileName($event)">
          </div>
          <button class="ui fluid large blue submit button" type="submit">Создать ревью</button>
        </div>
        <div class="ui negative message" v-if="error.length > 0">{{ error }}</div>
      </form>
    </div>
  </div>
</template>

<script>
export default {
  name: 'NewReview',
  data () {
    return {
      name: '',
      reviewers: '',
      filename: 'Добавьте файл',
      fileContent: '',
      error: ''
    }
  },
  methods: {
    changeFileName (event) {
      var file = event.target.files[0]
      if (file) {
        var reader = new FileReader()
        reader.readAsDataURL(file, 'UTF-8')
        var _this = this
        reader.onload = function (evt) {
          _this.filename = file.name
          _this.fileContent = evt.target.result.replace(/^data:.+;base64,/, '')
          _this.error = ''
        }
        reader.onerror = function (evt) {
          _this.filename = 'Добавьте файл'
          _this.fileContent = ''
          _this.error = 'Ошибка при чтении файла'
        }
      }
    },
    createNewReview () {
      if (this.fileContent.length === 0) {
        this.error = 'Добавьте файл для создания ревью'
        return
      }
      if (this.name.length === 0) {
        this.error = 'Имя необходимо'
        return
      }
      if (this.reviewers.length === 0) {
        this.error = 'Список ревьюеров необходим'
        return
      }
      this.$http.post('/reviews/new', {
        name: this.name,
        reviewers: this.reviewers,
        file_name: this.filename,
        file_content: this.fileContent
      }).then(() => {
        this.$router.push({name: 'OutReviews'})
      }).catch((err) => {
        if (!err.response.status) {
          this.error = 'Ошибка сети'
        } else if (err.response.status === 406) {
          this.error = 'Неверный список логинов ревьюеров'
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
