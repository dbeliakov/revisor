<template>
    <div class="comments">
        <div class="comment"
        v-for="comment in comments"
        :key="comment.id"
        :class="{'not-first': comment.id != comments[0].id}">
            <div class="avatar" :style="{'background-color': avatarColors[userHash(comment.author) % avatarColors.length]}">
                <span class="avatar-auto">{{comment.author.first_name[0]}}</span>
            </div>
            <div class="comment-body">
                <div class="header">
                    <span class="author">{{ comment.author.first_name }} {{ comment.author.last_name }}</span>
                    <span class="login">{{comment.author.username}}</span>
                </div>
                <div class="content">{{ comment.text }}</div>
                <div class="footer">
                    <a href="#" @click.prevent="$emit('new-comment')">Ответить</a><i class="circle icon"></i>
                    <!--<a href="#" @click.prevent>Редактировать</a><i class="circle icon"></i>-->
                    <a href="#" @click.prevent>{{ timeToString(new Date(comment.created * 1000)) }}</a>
                </div>
            </div>
        </div>
        <NewComment
        v-if="baseNewComment"
        class="new-comment"
        :class="{'not-first': comments.length > 0}"
        :author="$auth.user()"
        :reviewId="reviewId"
        :lineId="lineId"
        @saved="$emit('saved')"></NewComment>
    </div>
</template>

<script lang="ts">
import {Component, Vue, Prop, Provide} from 'vue-property-decorator';
import {timeToString} from '@/utils/utils';
import Comment from '@/reviews/comment';
import {UserInfo} from '@/auth/user-info';
import NewComment from '@/components/NewComment.vue';

@Component({
    components: {NewComment},
})
export default class Comments extends Vue {
    @Prop({default: []}) public readonly comments!: Comment[];
    @Prop({default: false}) public baseNewComment!: boolean;
    @Prop({default: ''}) public readonly reviewId!: string;
    @Prop({default: ''}) public readonly lineId!: string;

    public timeToString = timeToString;
    public avatarColors = [
        "#FFCC00", "#FF6666", "#CC66CC",
        "#9966FF", "#3366FF", "#66CCCC",
        "#33FF99", "#CCCC33", "#99CC33"]

    public userHash(user: UserInfo) {
        var hash = 0, i, chr;
        if (user.id.length === 0) return hash;
        for (i = 0; i < user.id.length; i++) {
            chr   = user.id.charCodeAt(i);
            hash  = ((hash << 5) - hash) + chr;
            hash |= 0; // Convert to 32bit integer
        }
        return hash;
    }
}
</script>

<style lang="scss">
.comments {
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
    white-space: pre-wrap;
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
