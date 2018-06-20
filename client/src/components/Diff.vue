<template>
  <div id="diff-content">
    <div class="d2h-wrapper">
          <div class="d2h-file-wrapper">
            <div class="d2h-file-diff">
              <div class="d2h-code-wrapper">
                <table class="d2h-diff-table">
                  <tbody class="d2h-diff-tbody">
                    <template v-for="group in groups">
                    <tr v-bind:key="group.id">
                      <td class="d2h-code-linenumber d2h-info"></td>
                      <td class="d2h-info">
                        <div class="d2h-code-line d2h-info">{{ group.range }}</div>
                      </td>
                    </tr>
                    <template v-for="line in group.lines" >
                      <tr v-bind:key="line.id" @click="showNewCommentForm(line.id)">
                        <td class="d2h-code-linenumber" v-bind:class="{'d2h-cntx': line.type === 'no', 'd2h-ins': line.type === 'insert', 'd2h-del': line.type === 'delete'}">
                          <div class="line-num1">{{ line.old_num }}</div>
                          <div class="line-num2">{{ line.new_num }}</div>
                        </td>
                        <td v-bind:class="{'d2h-cntx': line.type === 'no', 'd2h-ins': line.type === 'insert', 'd2h-del': line.type === 'delete'}">
                          <div class="d2h-code-line" v-bind:class="{'d2h-cntx': line.type === 'no', 'd2h-ins': line.type === 'insert', 'd2h-del': line.type === 'delete'}">
                            <span class="d2h-code-line-ctn">{{ line.old ? line.old.content : line.new.content }}</span>
                          </div>
                        </td>
                      </tr>
                      <tr v-bind:key="line.id + 'comment'" v-if="comments[line.id]">
                        <td colspan="2"><comments :comments="comments[line.id]" @new-comment="showNewCommentForm(line.id)"></comments></td>
                      </tr>
                      <tr v-bind:key="line.id + 'new_comment'" v-if="newCommentsShown[line.id]">
                        <td colspan="2"><new-comment :parent="''" :reviewId="reviewId" :lineId="line.id" @saved="$emit('update-all')"></new-comment></td>
                      </tr>
                    </template>
                    </template>
                  </tbody>
                </table>
              </div>
            </div>
          </div>
    </div>
  </div>
</template>

<script>
/* global Diff2HtmlUI */
import Comments from '@/components/Comments'
import NewComment from './NewComment.vue'

var $ = require('jquery')
window.$ = $
window.jQuery = $

require('highlightjs')
require('highlightjs/styles/github.css')
require('diff2html')
require('diff2html/dist/diff2html.css')
require('diff2html/dist/diff2html-ui.js')
require('semantic-ui-css/semantic.min.js')

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

export default {
  name: 'Diff',
  props: {
    diff: '',
    reviewId: '',
    commentsList: ''
  },
  components: {
    'comments': Comments,
    'new-comment': NewComment
  },
  data () {
    return {
      newCommentsShown: {}
    }
  },
  computed: {
    groups () {
      var result = []
      for (var i = 0; i < this.diff.groups.length; ++i) {
        var group = this.diff.groups[i]
        var oldFrom = group.old_range.from + 1
        var newFrom = group.new_range.from + 1
        var resGroup = {
          range: '@@ -' + formatRangeUnified(group.old_range.from, group.old_range.to) + ' +' +
            formatRangeUnified(group.new_range.from, group.new_range.to) + ' @@\n',
          lines: [],
          id: 'group' + i
        }
        for (var j = 0; j < group.lines.length; ++j) {
          var line = group.lines[j]
          if (line.type === 'no') {
            line.old_num = oldFrom++
            line.new_num = newFrom++
            if (line.old.id !== line.new.id) {
              console.log('Incorrect line ids')
            }
            line.id = line.old.id
          } else if (line.type === 'insert') {
            line.new_num = newFrom++
            line.id = line.new.id
          } else if (line.type === 'delete') {
            line.old_num = oldFrom++
            line.id = line.old.id
          }
          resGroup.lines.push(line)
        }
        result.push(resGroup)
      }
      return result
    },
    comments () {
      var result = {}
      for (var j = 0; j < this.commentsList.length; ++j) {
        if (!result[this.commentsList[j].line_id]) {
          result[this.commentsList[j].line_id] = []
        }
        // "Unexpected side effect in "comments" computed property" overvise
        var tmp = result[this.commentsList[j].line_id]
        tmp.push(this.commentsList[j])
        result[this.commentsList[j].line_id] = tmp
      }
      return result
    }
  },
  updated () {
    var diff2htmlUi = new Diff2HtmlUI()
    diff2htmlUi.highlightCode('#diff-content')
  },
  mounted () {
    this.newCommentsShown = {}
    var diff2htmlUi = new Diff2HtmlUI()
    diff2htmlUi.highlightCode('#diff-content')
  },
  methods: {
    showNewCommentForm (lineId) {
      if (this.newCommentsShown[lineId]) {
        this.newCommentsShown[lineId] = false
      } else {
        this.newCommentsShown[lineId] = true
      }
      this.$forceUpdate()
    }
  },
  watch: {
    diff () {
      this.newCommentsShown = {}
    }
  }
}
</script>
