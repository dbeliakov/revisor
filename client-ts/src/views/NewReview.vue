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
          <div class="field">
            <div class="ui input">
              <input name="reviewers" v-model="reviewers" placeholder="Логины ревьюеров (через запятую)" type="text">
            </div>
          </div>
          <!--<div class="field">
            <div class="ui input">
              <input name="reviewers" placeholder="Добавить ревьюера" type="text" @input="loadReviewersList">
            </div>
          </div>
          <div class=" dropdown-content" v-if="reviewersList.length > 0">
            <div class="item" v-for="reviewer in reviewersList" :key="reviewer">{{ reviewer }}</div>
          </div>-->
          <FileLoader @onStartReading='onStartReadingFile' @onFinishReading='onFinishReadingFile'/>
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

@Component({
    components: {FileLoader},
})
export default class NewReview extends Vue {
    public name: string = '';
    public reviewers: string = '';
    public filename: string = '';
    public fileContent: string = '';
    public error: string = '';
    public formDisabled: boolean = false;

    public onStartReadingFile() {
        this.disableForm();
    }

    public onFinishReadingFile(filename: string, content: string) {
        this.fileContent = content.replace(/^data:.+;base64,/, '');
        this.filename = filename;
        this.enableForm();
    }

    public async createReview() {
        if (this.formDisabled) {
            return;
        }
        this.disableForm();
        try {
            const result = await this.$reviews.createReview(
              this.name, this.reviewers, this.filename, this.fileContent);
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
</style>

