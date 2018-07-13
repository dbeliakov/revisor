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
              <input name="reviewers" v-model="reviewers" placeholder="Логины ревьюеров (через запятую)" type="text">
            </div>
          </div>
          <div class="field">
            <label for="file" class="ui icon button">
              <i class="file icon"></i>
              {{ filename }}</label>
            <input type="file" id="file" style="display:none" v-on:change="updateFile($event)">
          </div>
          <button class="ui fluid large blue submit button" v-bind:class="{'disabled': formDisabled}" type="submit">Создать ревью</button>
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
      error: '',
      formDisabled: false
    }
  },
  methods: {
    updateFile (event) {
      if (this.formDisabled) {
        return
      }
      var file = event.target.files[0]
      if (file) {
        if (file.size > 1048576) { // 1mb
          this.error = 'Максимальный размер файла: 1 мегабайт'
          return
        }
        this.formDisabled = true
        var reader = new FileReader()
        var _this = this
        reader.onload = function (evt) {
          _this.filename = file.name
          _this.fileContent = evt.target.result.replace(/^data:.+;base64,/, '')
          _this.error = ''
          _this.formDisabled = false
        }
        reader.onerror = function (evt) {
          _this.filename = 'Добавьте файл'
          _this.fileContent = ''
          _this.error = 'Ошибка при чтении файла'
          _this.formDisabled = false
        }
        reader.readAsDataURL(file, 'UTF-8')
      }
    },
    createNewReview () {
      if (this.formDisabled) {
        return
      }
      this.formDisabled = true
      if (this.fileContent.length === 0) {
        this.error = 'Добавьте файл для создания ревью'
        this.formDisabled = false
        return
      }
      if (this.name.length === 0) {
        this.error = 'Имя необходимо'
        this.formDisabled = false
        return
      }
      if (this.reviewers.length === 0) {
        this.error = 'Список ревьюеров необходим'
        this.formDisabled = false
        return
      }
      this.$http.post('/reviews/new', {
        name: this.name,
        reviewers: this.reviewers,
        file_name: this.filename,
        file_content: this.fileContent
      }).then(() => {
        this.formDisabled = false
        this.$router.push({name: 'OutReviews'})
      }).catch((err) => {
        if (!err.response.status) {
          this.error = 'Ошибка сети'
          this.formDisabled = false
          return
        }
        if (err.response.data.client_message) {
          this.error = err.response.data.client_message
        } else {
          this.error = 'Внутренняя ошибка сервиса'
        }
        this.formDisabled = false
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
