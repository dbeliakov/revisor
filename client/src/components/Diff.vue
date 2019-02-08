<template>
  <div id="diff-content">
    <div class="d2h-wrapper">
          <div class="d2h-file-wrapper" :data-lang="fileExt()">
            <div class="d2h-file-diff">
              <div class="d2h-code-wrapper">
                <table class="d2h-diff-table">
                  <tbody class="d2h-diff-tbody">
                    <template v-for="group in computedGroups()">
                    <template v-for="line in group.lines" >
                      <tr v-bind:key="line.id" @click="showNewCommentForm(line.id)">
                        <td class="d2h-code-linenumber" v-bind:class="{'d2h-cntx': line.type === 'no', 'd2h-ins': line.type === 'insert', 'd2h-del': line.type === 'delete'}">
                          <div class="line-num1" v-if="line.oldNum > 0">{{ line.oldNum }}</div>
                          <div class="line-num2" v-if="line.newNum > 0">{{ line.newNum }}</div>
                        </td>
                        <td v-bind:class="{'d2h-cntx': line.type === 'no', 'd2h-ins': line.type === 'insert', 'd2h-del': line.type === 'delete'}">
                          <div class="d2h-code-line" v-bind:class="{'d2h-cntx': line.type === 'no', 'd2h-ins': line.type === 'insert', 'd2h-del': line.type === 'delete'}">
                            <span class="d2h-code-line-ctn" v-html="line.old ? line.hOldLine : line.hNewLine"></span>
                          </div>
                        </td>
                      </tr>
                      <tr
                      v-bind:key="line.id + 'comment'"
                      v-if="computedComments()[line.id] || newCommentsShown[line.id]"
                      v-bind:class="{'d2h-cntx': line.type === 'no', 'd2h-ins': line.type === 'insert', 'd2h-del': line.type === 'delete'}">
                        <td colspan="2"><Comments
                          :newCommentFormShown="newCommentsShown[line.id]"
                          :comments="computedComments()[line.id]"
                          :reviewId="reviewId"
                          :lineId="line.id"
                          @saved="$emit('update-all')"
                          @cancelled="newCommentsShown[line.id] = false; $forceUpdate()"></Comments></td>
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

<script lang="ts">
import {Component, Vue, Prop, Watch} from 'vue-property-decorator';
import { Diff, DiffLine, Line } from '@/reviews/diff';
import Comment from '@/reviews/comment';
import Comments from '@/components/Comments.vue';
import hljs from 'highlightjs';

class UILine {
  public oldNum: number = 0;
  public newNum: number = 0;
  public id: string = '';

  public old: Line | undefined;
  public new: Line | undefined;
  public type: string;

  public hOldLine: string | undefined;
  public hNewLine: string | undefined;

  constructor(line: DiffLine) {
    this.old = line.old;
    this.new = line.new;
    this.type = line.type;
  }
}

@Component({
  components: {Comments},
})
export default class DiffComponent extends Vue {
  @Prop({default: undefined}) public diff!: Diff;
  @Prop({default: ''}) public reviewId!: number;
  @Prop({default: []}) public commentsList!: Comment[];

  public newCommentsShown: {[key: string]: boolean} = {};

  public computedGroups() {
      const result = [];
      let oldContinuation;
      let newContinuation;
      for (let i = 0; i < this.diff.groups.length; ++i) {
        const group = this.diff.groups[i];
        let oldFrom = group.oldRange.from + 1;
        let newFrom = group.newRange.from + 1;
        const resGroup = {
          id: 'group' + i,
          lines: new Array<UILine>(),
        };
        for (const line of group.lines) {
          const uiLine = new UILine(line);
          if (uiLine.old) {
            const hLine = hljs.highlight(this.fileExt(), uiLine.old.content, true, oldContinuation);
            uiLine.hOldLine = hLine.value;
            oldContinuation = hLine.top;
          }
          if (uiLine.new) {
            const hLine = hljs.highlight(this.fileExt(), uiLine.new.content, true, newContinuation);
            uiLine.hNewLine = hLine.value;
            newContinuation = hLine.top;
          }
          if (line.type === 'no') {
            uiLine.oldNum = oldFrom++;
            uiLine.newNum = newFrom++;
            uiLine.id = line.old!.id; // === line.new.id
          } else if (line.type === 'insert') {
            uiLine.newNum = newFrom++;
            uiLine.id = line.new!.id;
          } else if (line.type === 'delete') {
            uiLine.oldNum = oldFrom++;
            uiLine.id = line.old!.id;
          }
          resGroup.lines.push(uiLine);
        }
        result.push(resGroup);
      }
      return result;
    }

    public computedComments() {
      const result: {[key: string]: any} = {};
      for (const comment of this.commentsList) {
        if (!result[comment.lineId]) {
          result[comment.lineId] = [];
        }
        // "Unexpected side effect in "comments" computed property overvise
        // TODO
        const tmp = result[comment.lineId];
        tmp.push(comment);
        result[comment.lineId] = tmp;
      }
      return result;
    }

    public fileExt(): string {
        return this.diff.filename.split('.').pop() as string;
    }

    public mounted() {
      this.newCommentsShown = {};
    }

    public showNewCommentForm(lineId: string) {
      this.newCommentsShown[lineId] = true;
      this.$forceUpdate();
    }

    @Watch('diff')
    public onDiffChanged() {
      this.newCommentsShown = {};
    }
}
</script>

<style lang="scss" scoped>
@import '~highlightjs/styles/github.css';
@import '~diff2html/dist/diff2html.css';

.d2h-file-diff {
  overflow-x: hidden !important;
}
</style>
