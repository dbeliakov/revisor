<template>
  <div id="new_review">
    <h1>Новое ревью</h1>
    <div class="column">
      <form class="ui large form" v-on:submit.prevent="createReview()">
        <div class="ui">
          <div class="field">
            <div class="ui input">
              <input name="name" v-model="name" placeholder="Заголовок" type="text">
            </div>
          </div>
          <div class="reviewers-dropdown">
            <div class="field" style="margin-bottom: 0px !important;">
                <div class="ui input">
                    <input
                        name="reviewers"
                        v-model="reviewerSearch"
                        placeholder="Добавить ревьюера"
                        type="text"
                        autocomplete="off"
                        @input="loadReviewersList"
                        @keydown.enter.prevent>
                </div>
            </div>
            <div class="ui vertical menu reviewers-list" style="width: 100%;" v-if="reviewersList.length > 0">
                <a class="item"
                    v-for="reviewer in reviewersList"
                    :key="reviewer.id"
                    @click="addReviewer(reviewer.username)">
                    {{reviewer.firstName}} {{reviewer.lastName}} ({{reviewer.username}})
                </a>
            </div>
          </div>
          <div class="field" style="text-align: left; font-size: 14pt; margin-top: 20px;">
                <b>Ревьюеры: </b><span v-if="reviewers.length === 0">
                  нет добавленных ревьюеров
                </span><span v-for="reviewer in reviewers" :key="reviewer.id">{{reviewer}}
                    <i style="font-size: 12pt; cursor: pointer;"
                        @click="removeReviewer(reviewer)"
                        class="window close outline icon"></i></span>
          </div>
          <FileLoader
            @onStartReading='onStartReadingFile'
            @onFinishReading='onFinishReadingFile'
            @onReadingError='onReadingError'/>
          <button class="ui fluid large blue submit button" v-bind:class="{'disabled': formDisabled}" type="submit">Создать ревью</button>
        </div>
        <div class="ui negative message" v-if="error.length > 0">{{ error }}</div>
      </form>
    </div>
  </div>
</template>

<script lang="ts">
import {Component, Vue} from 'vue-property-decorator';
import FileLoader from '@/components/FileLoader.vue';
import { UserInfo } from '@/auth/user-info';
import { error } from 'util';

@Component({
    components: {FileLoader},
})
export default class NewReview extends Vue {
    public name: string = '';
    public reviewers: string[] = [];
    public filename: string = '';
    public fileContent: string = '';
    public error: string = '';
    public formDisabled: boolean = false;
    public reviewersList: UserInfo[] = [];
    public reviewerSearch: string = '';

    public onStartReadingFile() {
        this.disableForm();
    }

    public onFinishReadingFile(filename: string, content: string) {
        this.fileContent = content.replace(/^data:.+;base64,/, '');
        this.filename = filename;
        this.error = '';
        this.enableForm();
    }

    public onReadingError(err: string) {
        this.error = err;
        this.filename = '';
        this.fileContent = '';
        this.enableForm();
    }

    public async loadReviewersList() {
        if (this.reviewerSearch.length === 0) {
            this.reviewersList = [];
            return;
        }
        const result = await this.$reviews.searchReviewers(this.reviewerSearch);
        if (result instanceof Error) {
            alert(result);
            return;
        }
        this.reviewersList = result;
    }

    public addReviewer(username: string) {
        if (this.reviewers.indexOf(username) === -1) {
            this.reviewers.push(username);
        }
        this.reviewersList = [];
        this.reviewerSearch = '';
    }

    public removeReviewer(username: string) {
        const index = this.reviewers.indexOf(username);
        if (index > -1) {
            this.reviewers.splice(index, 1);
        }
    }

    public async createReview() {
        if (this.formDisabled) {
            return;
        }
        this.disableForm();
        try {
            if (!this.validateForm()) {
                return;
            }
            const result = await this.$reviews.createReview(
              this.name, this.reviewers.join(','), this.filename, this.fileContent);
            if (result) {
              this.error = result.message;
            } else {
              this.$router.push({
                name: 'outgoing',
              });
            }
        } finally {
            this.enableForm();
        }
    }

    private validateForm(): boolean {
        if (this.name.length === 0) {
            this.error = 'Введите заголовок ревью';
            return false;
        }
        if (this.reviewers.length === 0) {
            this.error = 'Добавьте как минимум одного ревьюера';
            return false;
        }
        if (this.filename.length === 0 || this.fileContent.length === 0) {
            this.error = 'Добавьте файл';
            return false;
        }
        return true;
    }

    private disableForm() {
        this.formDisabled = true;
    }

    private enableForm() {
        this.formDisabled = false;
    }
}
</script>

<style lang="scss">
@import '~semantic-ui-css/semantic.min.css';

body > .grid {
    height: 100%;
}
.column {
    max-width: 450px;
    margin: auto;
    margin-top: 40px;
}

.reviewers-list {
    margin-top: 0px !important;
    border-top-left-radius: 0px;
    border-top-right-radius: 0px;
}
</style>
