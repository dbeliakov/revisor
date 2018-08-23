<template>
    <div class="line-comments">
        <CommentComponent
        v-for="comment in comments"
        :key="comment.id"
        :notFirst="comment.id != comments[0].id"
        :comment="comment"
        :reviewId="reviewId"
        :lineId="lineId"
        @saved="$emit('saved')"></CommentComponent>

        <NewComment
        v-if="baseNewComment"
        class="new-comment"
        :class="{'not-first': comments.length > 0}"
        :author="$auth.user()"
        :reviewId="reviewId"
        :lineId="lineId"
        @saved="$emit('saved')"
        @cancelled="$emit('cancelled')"></NewComment>
    </div>
</template>

<script lang="ts">
import {Component, Vue, Prop, Provide} from 'vue-property-decorator';
import {timeToString} from '@/utils/utils';
import Comment from '@/reviews/comment';
import {UserInfo} from '@/auth/user-info';
import NewComment from '@/components/NewComment.vue';
import CommentComponent from '@/components/Comment.vue';
import Marked from 'marked';

@Component({
    components: {NewComment, CommentComponent},
})
export default class Comments extends Vue {
    @Prop({default: []}) public readonly comments!: Comment[];
    @Prop({default: false}) public baseNewComment!: boolean;
    @Prop({default: ''}) public readonly reviewId!: string;
    @Prop({default: ''}) public readonly lineId!: string;
}
</script>

<style lang="scss">
.line-comments {
    margin: 10px;
    border: 1px solid lightgrey;
    border-radius: 5px;
    background-color: white;
}

.not-first {
    border-top: 1px solid lightgrey;
    padding-top: 10px;
}

.comment,.new-comment {
    margin: 20px;
    margin-bottom: 10px;
    margin-top: 10px;
    display: flex;
}

.avatar {
    background-color: #f84848;
    color: #fff;
    border-radius: 50%;
    font-size: 25px;
    width: 40px;
    height: 40px;
    line-height: 40px;
    text-align: center;
    vertical-align: middle;
    border: 0.5px solid darkgray;
    margin-right: 10px;
    flex: none;
}

.content {
    margin-bottom: 5px;
    margin-top: 5px;
}

.author {
    color: #003399;
    font-weight: bold;
    font-size: 11pt;
}

.login {
    color: grey;
    font-size: 9pt;
    margin-left: 5px;
}

.footer {
    font-size: 9pt;
    font-weight: bold;
    a {
        color: grey;
    }
    .icon {
        color: grey;
        font-size: 3pt;
        vertical-align: top;
        margin-left: 5px;
        margin-right: 5px;
    }
    a:hover {
        text-decoration: underline;
    }
}

.comment-body {
    flex: 2;
}
</style>
