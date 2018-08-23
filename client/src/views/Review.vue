<template>
  <div style="margin: 20px 40px; text-align: left;">
    <div v-if="review" style="">
      <h1>{{ review.info.name }}</h1>
      <div v-if="!review.info.closed" style="margin-bottom: 15px;">
        <button v-if="review.info.owner.id === $auth.user().id" class="review-button ui primary basic button" @click="openModal">Обновить</button>
        <button v-if="review.info.owner.id !== $auth.user().id" class="review-button ui positive basic button" @click="accept()">Принять</button>
        <button class="review-button ui negative basic button" @click="decline()">Отклонить</button><br>
      </div>
      <h4 style="display: inline;">Создатель:</h4> {{ review.info.owner.first_name }}
        {{ review.info.owner.last_name }} ({{ review.info.owner.username }})<br>
      <h4 style="display: inline;">Ревьюеры:</h4>
        <span v-for="reviewer in review.info.reviewers" :key="reviewer.username">
            {{reviewer.first_name}} {{reviewer.last_name}} ({{reviewer.username}})
            <template v-if="reviewer !== review.info.reviewers[review.info.reviewers.length - 1]">,</template>
        </span><br>
      <template v-if="review.info.closed"><h4 style="display: inline;">Закрыто:</h4> <span v-if="review.info.accepted">Принято</span> <span v-if="!review.info.accepted">Отклонено</span><br></template>
      <h4 style="margin-bottom: 20px; display: inline;">Обновлено:</h4> {{timeConverter(review.info.updated)}}
    </div>
    <div v-if="review" style="text-align: center; margin-bottom: 20px;">
      <h4 style="margin-top: 30px; margin-bottom: 10px;" v-if="review.info.revisions_count > 1 && startRev !== null && endRev !== null">
        <span v-if="startRev != endRev">revision {{startRev}} <i class="right arrow icon"></i> revision {{endRev}}</span>
        <span v-else>revision {{startRev}}</span>
      </h4>
      <div v-if="review.info.revisions_count > 1" id="slider_wrapper" style="margin: 0px auto; width:300px;">
        <div id="slider_revisions"></div>
      </div>
    </div>

    <diff v-if="review" :diff="new Diff(review.diff)" :commentsList="review.comments.map((json) => new Comment(json))" :reviewId="review_id" @update-all="updateReview()"></diff>

    <div class="ui modal" id="add_revision">
      <i class="close icon"></i>
      <div class="header">
        Обновить ревью
      </div>
      <div class="content">
        <form class="ui large form" v-on:submit.prevent="updateReviewOnServer()">
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
            <button class="ui fluid large blue submit button" type="submit">Обновить ревью</button>
          </div>
          <div class="ui negative message" v-if="error.length > 0">{{ error }}</div>
        </form>
      </div>
    </div>
  </div>
</template>

<script>
/* tslint:disable */
import DiffComponent from '@/components/Diff'
import { Diff } from '@/reviews/diff';
import Comment from '@/reviews/comment';

var $ = require('jquery')
window.$ = $
window.jQuery = $

require('semantic-ui-css/semantic.min.js')
require('jquery-ui/ui/widgets/slider.js')
require('jquery-ui/themes/base/all.css')

export default {
  name: 'Review',
  components: {
    'diff': DiffComponent
  },
  props: ['id'],
  created () {
    this.updateReview()
  },
  data () {
    return {
      review_id: null,
      review: null,
      name: '',
      reviewers: '',
      error: '',
      filename: 'Добавить файл',
      fileContent: '',
      startRev: null,
      endRev: null,
      Diff: Diff,
      Comment: Comment,
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
    updateReview () {
      this.review_id = this.$route.params.id
      var params = {}
      if (this.startRev !== null) {
        params['start_rev'] = this.startRev - 1
      }
      if (this.endRev !== null) {
        params['end_rev'] = this.endRev - 1
      }
      this.$auth.axios.get('/reviews/' + this.review_id, {
        params: params
      }).then((response) => {
        this.review = response.data.data
        this.name = this.review.info.name
        this.reviewers = ''
        for (var i = 0; i < this.review.info.reviewers.length; ++i) {
          this.reviewers += this.review.info.reviewers[i].username
          if (i !== this.review.info.reviewers.length - 1) {
            this.reviewers += ','
          }
        }

        if (this.startRev === null) {
          this.startRev = 1
        }
        if (this.endRev === null) {
          this.endRev = this.review.info.revisions_count
        }

        if (this.review.info.revisions_count > 1) {
          var _this = this
          $(function () {
            $('#slider_revisions').slider({
              range: true,
              min: 1,
              max: _this.review.info.revisions_count,
              values: [_this.startRev, _this.endRev],
              slide: function (event, ui) {
                _this.startRev = ui.values[0]
                _this.endRev = ui.values[1]
                _this.updateReview()
              }
            })
            var width = Math.min(800, _this.review.info.revisions_count * 40)
            $('#slider_wrapper').css('width', '' + width + 'px')
          })
        }
      })
    },
    updateReviewOnServer () {
      var data = {
        name: this.name,
        reviewers: this.reviewers
      }
      if (this.fileContent.length > 0) {
        data['new_revision'] = this.fileContent
      }
      this.$auth.axios.post('/reviews/' + this.review_id + '/update', data).then(() => {
        $('#add_revision').modal('hide')
        this.startRev = null
        this.endRev = null
        this.filename = 'Добавьте файл'
        this.fileContent = ''
        this.updateReview()
      })
    },
    timeConverter (timestamp) {
      var toStr = function (val) {
        if (val < 10) {
          return '0' + val
        }
        return '' + val
      }

      var a = new Date(timestamp * 1000)
      var months = ['Января', 'Февраля', 'Марта', 'Апреля', 'Мая', 'Июня', 'Июля', 'Августа',
        'Сентября', 'Октября', 'Ноября', 'Декабря']
      var year = a.getFullYear()
      var month = months[a.getMonth()]
      var date = a.getDate()
      var hour = a.getHours()
      var min = a.getMinutes()
      var time = date + ' ' + month + ' ' + year + ' ' + toStr(hour) + ':' + toStr(min)
      return time
    },
    openModal () {
      $('#add_revision').modal('show')
    },
    accept () {
      this.$auth.axios.get('/reviews/' + this.review_id + '/accept').then(() => {
        this.updateReview()
      })
    },
    decline () {
      this.$auth.axios.get('/reviews/' + this.review_id + '/decline').then(() => {
        this.updateReview()
      })
    }
  },
  watch: {
    $route (to, from) {
      this.review = {}
      this.updateReview()
    }
  }
}
</script>

<style scoped>
h4 {
  margin: 0px;
}

.d2h-code-line del,
.d2h-code-side-line del {
  background-color: #fee8e9 !important;
}

.d2h-code-line ins,
.d2h-code-side-line ins {
  background-color: #dfd;
}

.d2h-code-linenumber {
  width: 66px;
}

.line-num1 {
  width: 30px;
}

.line-num2 {
  width: 30px;
}

.review-button {
  width: 100px;
  padding: 5px !important;
}
</style>
