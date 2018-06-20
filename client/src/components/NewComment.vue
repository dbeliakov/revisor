<template>
  <div>
    <div class="ui form comment">
      <div class="field">
        <textarea rows="2" class="text" v-model="text"></textarea>
        <div class="ui submit button primary" style="float: right;" @click="submit()">Отправить</div>
      </div>
    </div>
  </div>
</template>

<script>
export default {
  name: 'NewComment',
  props: ['parent', 'reviewId', 'lineId'], // TODO use parent
  data () {
    return {
      text: '',
      error: ''
    }
  },
  methods: {
    submit () {
      this.$http.post('/comments/add', {
        'line_id': this.lineId,
        'review_id': this.reviewId,
        'text': this.text,
        'parent': ''
      }).then(() => {
        this.$emit('saved')
      }).catch((err) => {
        if (!err.response.status) {
          this.error = 'Ошибка сети'
        } else if (err.response.status === 406) {
          this.error = 'Неверный список логинов ревьюеров'
        } else {
          this.error = 'Неизвестная ошибка'
        }
        if (this.error !== '') {
          alert(this.error)
        }
      })
    }
  }
}
</script>

<style scoped>
.comment {
  margin-left: 64px;
  margin-bottom: 5px;
}

.text {
  margin-bottom: 10px !important;
}

.header {
  margin-bottom: 10px;
}

.footer {
  margin-top: 10px;
}
</style>
