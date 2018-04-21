<template>
  <div style="margin: 20px 40px;">
    <div v-if="review" style="margin-bottom: 40px;">
      <h1>{{ review.info.name }}</h1>
      <h4>Создатель: {{ review.info.owner.first_name }}
        {{ review.info.owner.last_name }} ({{ review.info.owner.username }})</h4>
      <h4>Ревьюеры:
        <span v-for="reviewer in review.info.reviewers" :key="reviewer.username">
              {{reviewer.first_name}} {{reviewer.last_name}} ({{reviewer.username}})
              <span v-if="reviewer !== review.info.reviewers[review.info.reviewers.length - 1]">,</span>
            </span>
      </h4>
      <h4 v-if="review.info.closed">Закрыто: <span v-if="review.info.accepted">Принято</span> <span v-if="!review.info.accepted">Отклонено</span></h4>
      <h4 style="margin-bottom: 20px;">Обновлено: {{timeConverter(review.info.updated)}}</h4>

      <div v-if="!review.info.closed">
      <button v-if="review.info.owner.id === $auth.user().id" class="ui primary basic button" @click="openModal">Обновить ревью</button>
      <button v-if="review.info.owner.id !== $auth.user().id" class="ui positive basic button" @click="accept()">Принять</button>
      <button class="ui negative basic button" @click="decline()">Отклонить</button><br>
      </div>

      <h4 style="margin-top: 30px; margin-bottom: 10px;" v-if="review.info.revisions_count > 1 && startRev !== null && endRev !== null">
        rev{{startRev}} -> rev{{endRev}}
      </h4>
      <div v-if="review.info.revisions_count > 1" id="slider_wrapper" style="margin: 0px auto; width:300px;">
        <div id="slider_revisions"></div>
      </div>
    </div>

    <div id="diff_view"></div>

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
/* global Diff2HtmlUI */
var $ = require('jquery')
window.$ = $
window.jQuery = $

require('highlightjs')
require('highlightjs/styles/github.css')
require('diff2html')
require('diff2html/dist/diff2html.css')
require('diff2html/dist/diff2html-ui.js')
require('semantic-ui-css/semantic.min.js')
require('jquery-ui/ui/widgets/slider.js')
require('jquery-ui/themes/base/all.css')

function formatRangeUnified (start, stop) {
  var beginning = start + 1
  var length = stop - start
  if (length === 1) {
    return '' + beginning
  }
  if (length === 0) {
    --beginning
  }
  return '' + beginning + ',' + length
}

function toDiffString (content) {
  var result = ''
  result += '--- ' + content.filename + '\n'
  result += '+++ ' + content.filename + '\n'

  for (var i = 0; i < content.groups.length; ++i) {
    var group = content.groups[i]
    result += '@@ -' + formatRangeUnified(group.old_range.from, group.old_range.to) + ' +' +
      formatRangeUnified(group.new_range.from, group.new_range.to) + ' @@\n'
    for (var l = 0; l < group.lines.length; ++l) {
      var line = group.lines[l]
      if (line.type === 'no') {
        result += ' ' + line.old.content + '\n'
      } else if (line.type === 'insert') {
        result += '+' + line.new.content + '\n'
      } else if (line.type === 'delete') {
        result += '-' + line.old.content + '\n'
      }
    }
  }

  return result
}

export default {
  name: 'Review',
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
      endRev: null
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
      this.$http.get('/reviews/' + this.review_id, {
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
        console.log(data['new_revision'])
      } else {
        console.log('NO NEW REVISION')
      }
      this.$http.post('/reviews/' + this.review_id + '/update', data).then(() => {
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
      this.$http.get('/reviews/' + this.review_id + '/accept').then(() => {
        this.updateReview()
      })
    },
    decline () {
      this.$http.get('/reviews/' + this.review_id + '/decline').then(() => {
        this.updateReview()
      })
    }
  },
  watch: {
    $route (to, from) {
      this.review = {}
      this.updateReview()
    },
    review () {
      var diff2htmlUi = new Diff2HtmlUI({diff: toDiffString(this.review.diff)})
      diff2htmlUi.draw('#diff_view', {
        inputFormat: 'json',
        outputFormat: 'line-by-line',
        showFiles: false,
        matching: 'none',
        rawTemplates: {
          'generic-line': `<tr>
    <td class="{{lineClass}} {{type}}">
      {{{lineNumber}}}
    </td>
    <td class="{{type}}">
        <div class="{{contentClass}} {{type}}">
        {{#content}}
            <span class="d2h-code-line-ctn">{{{content}}}</span>
        {{/content}}
        </div>
    </td>
</tr>`
        }
      })
      diff2htmlUi.highlightCode('#diff_view')

      console.log($('tbody').length)
      $('.d2h-diff-tbody').find('tr').each(function () {
        var el = $(this)
        el.find('.d2h-code-linenumber').click(() => {
          $('<tr><td class="d2h-code-linenumber d2h-info"></td><td class="d2h-cntx">Комментарий</td></tr>').insertAfter(el)
        })
      })
    }
  }
}
</script>

<!-- Add "scoped" attribute to limit CSS to this component only -->
<style>
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
</style>
