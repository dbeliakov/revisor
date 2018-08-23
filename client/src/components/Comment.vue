<template>
    <div style="display: inline;">
        <div class="comment" :class="{'not-first': notFirst}" :style="{'margin-left': (20 + level * 40) + 'px'}">
            <div class="avatar" :style="{'background-color': userAvatarColor(comment.author)}">
                <span class="avatar-auto">{{comment.author.firstName[0]}}</span>
            </div>
            <div class="comment-body">
                <div class="header">
                    <span class="author">{{ comment.author.firstName }} {{ comment.author.lastName }}</span>
                    <span class="login">{{comment.author.username}}</span>
                </div>
                <div class="content"><span v-html="toMarkdown(comment.text)"></span></div>
                <div class="footer">
                    <a href="#" @click.prevent="showReply=true;">Ответить</a><i class="circle icon"></i>
                    <a href="#" @click.prevent>{{ timeToString(comment.created) }}</a>
                </div>
            </div>
        </div>

        <Comment
        v-for="child in comment.childs"
        :key="child.id"
        class='not-first'
        :level="level + 1"
        :reviewId="reviewId"
        :lineId="lineId"
        :comment="child"
        @saved="$emit('saved')"></Comment>

        <NewComment
        :style="{'margin-left': (20 + (1 + level) * 40) + 'px'}"
        v-if="showReply"
        class="new-comment not-first"
        :author="$auth.user()"
        :reviewId="reviewId"
        :lineId="lineId"
        :parentId="comment.id"
        @saved="$emit('saved'); showReply=false;"
        @cancelled="showReply = false;"></NewComment>
    </div>
</template>


<script lang="ts">
import {Component, Vue, Prop, Provide} from 'vue-property-decorator';
import { timeToString, userAvatarColor } from '@/utils/utils';
import { UserInfo } from '@/auth/user-info';
import Marked from 'marked';
import NewComment from '@/components/NewComment.vue';

@Component({
    components: {NewComment},
} )
export default class Comment extends Vue {
    @Prop({default: undefined}) public readonly comment!: Comment;
    @Prop({default: true}) public notFirst!: boolean;
    @Prop({default: ''}) public readonly reviewId!: string;
    @Prop({default: ''}) public readonly lineId!: string;
    @Prop({default: 0}) public readonly level!: number;

    public showReply: boolean = false;

    public timeToString = timeToString;
    public userAvatarColor = userAvatarColor;

    public toMarkdown(text: string) {
        return Marked.parse(text);
    }
}
</script>
