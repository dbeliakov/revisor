<template>
  <div style="margin: 20px 40px;">
    <div class="ui checkbox" style="float: left; margin-bottom: 10px;">
      <input name="example" type="checkbox" v-model="showClosed">
      <label>Отображать закрытые</label>
    </div>
    <table class="ui celled striped table">
      <thead>
        <tr><th>
          Имя
        </th>
        <th>
          Создатель
        </th>
        <th>
          Ревьюеры
        </th>
        <th>
          Обновлено
        </th>
      </tr></thead>
      <tbody>
        <tr v-for="review in computedFilteredReviews" :key="review.id"
            v-bind:class="{positive: review.closed && review.accepted, negative: review.closed && !review.accepted}">
          <td><router-link :to="'/review/' + review.id">{{review.name}}</router-link></td>
          <td class="collapsing">
            {{review.owner.firstName}} {{review.owner.lastName}} ({{review.owner.username}})
          </td>
          <td class="collapsing">
            <span v-for="reviewer in review.reviewers" :key="reviewer.username">
              {{reviewer.firstName}} {{reviewer.lastName}} ({{reviewer.username}})
            </span>
          </td>
          <td class="right aligned collapsing">{{timeToString(review.updated)}}</td>
        </tr>
      </tbody>
    </table>
  </div>
</template>

<script lang="ts">
import {Component, Vue, Prop, Watch} from 'vue-property-decorator';
import Review from '@/reviews/review';
import ReviewsService from '@/reviews/service';
import {timeToString} from '@/utils/utils';

@Component
export default class Reviews extends Vue {
    @Prop(String)
    public type!: string;

    public reviews: Review[] = [];
    public error: string = '';
    public showClosed: boolean = false;

    public timeToString = timeToString;
    private service?: ReviewsService = undefined;

    public created() {
        this.service = new ReviewsService(this.$auth.axios); // TODO fix it
        this.updateReviews();
    }

    public async updateReviews() {
        let result: Review[] | Error;
        if (this.type === 'incoming') {
            result = await this.service!.loadIncomingReviews();
        } else if (this.type === 'outgoing') {
            result = await this.service!.loadOutgoingReviews();
        } else {
            return;
        }
        if (result instanceof Error) {
            this.error = result.message;
            return;
        }
        this.reviews = result;
        this.error = '';
    }

    public get computedFilteredReviews() {
        if (this.showClosed) {
            return this.reviews;
        }
        return this.reviews.filter((review) => !review.closed);
    }

    @Watch('$route')
    public onRouteChanged(from: any, to: any) {
        this.reviews = [];
        this.updateReviews();
    }
}
</script>