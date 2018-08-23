<template>
    <div class="new-comment">
      <div class="avatar" :style="{'background-color': userAvatarColor(author)}">
          <span class="avatar-auto">{{author.firstName[0]}}</span>
      </div>
      <div class="comment-body">
          <div class="header">
              <span class="author">{{ author.firstName }} {{ author.lastName }}</span>
              <span class="login">{{author.username}}</span>
          </div>
          <div class="content ui form">
              <textarea v-if="!previewEnabled" rows="4" class="text" v-model="text"></textarea>
              <span v-if="previewEnabled" v-html="computedMerkedText()"></span>
          </div>
          <div class="footer">
              <a href="#" @click.prevent="submit()">Добавить</a><i class="circle icon"></i>
              <template v-if="!previewEnabled"><a href="#" @click.prevent="previewEnabled = true;">Предпросмотр</a><i class="circle icon"></i></template>
              <template v-if="previewEnabled"><a href="#" @click.prevent="previewEnabled = false;">Редактировать</a><i class="circle icon"></i></template>
              <a href="#" @click.prevent="$emit('cancelled')">Отменить</a>
          </div>
      </div>
    </div>
</template>

<script lang="ts">
import {Component, Vue, Prop} from 'vue-property-decorator';
import Comment from '@/reviews/comment';
import { UserInfo } from '@/auth/user-info';
import Marked from 'marked';
import { userAvatarColor } from '@/utils/utils';

@Component
export default class NewComment extends Vue {
    @Prop({default: ''}) public readonly reviewId!: string;
    @Prop({default: ''}) public readonly lineId!: string;
    @Prop({default: undefined}) public readonly author!: UserInfo;
    @Prop({default: ''}) public readonly parentId!: string;

    public text: string = '';
    public previewEnabled: boolean = false;

    public userAvatarColor = userAvatarColor;

    public async submit() {
      const error = await this.$reviews.addComment(this.lineId, this.reviewId, this.text, this.parentId);
      if (error) {
        alert(error.message);
      } else {
        this.$emit('saved');
      }
    }

    public computedMerkedText() {
        return Marked.parse(this.text);
    }
}
</script>

<style scoped>
.text {
  width: 100%;
}
</style>
