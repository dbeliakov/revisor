<template>
  <div style="margin: 20px 40px;">
    <table class="ui celled striped table">
      <thead>
        <tr><th>
          Имя
        </th>
        <th>
          Создатель
        </th>
        <th>
          Проверяющие
        </th>
        <th>
          Обновлено
        </th>
      </tr></thead>
      <tbody>
        <tr v-for="review in reviews" :key="review.id"
            v-bind:class="{positive: review.closed && review.accepted, negative: review.closed && !review.accepted}">
          <td><router-link :to="'/review/' + review.id">{{review.name}}</router-link></td>
          <td class="collapsing">
            {{review.owner.first_name}} {{review.owner.last_name}} ({{review.owner.username}})
          </td>
          <td class="collapsing">
            <span v-for="reviewer in review.reviewers" :key="reviewer.username">
              {{reviewer.first_name}} {{reviewer.last_name}} ({{reviewer.username}})
            </span>
          </td>
          <td class="right aligned collapsing">{{timeConverter(review.updated)}}</td>
        </tr>
      </tbody>
    </table>
  </div>
</template>

<script>
export default {
  name: 'Reviews',
  props: ['type'],
  created () {
    this.$http.get('/reviews/' + this.type).then((response) => {
      this.reviews = response.data.data
    })
  },
  data () {
    return {
      reviews: []
    }
  },
  methods: {
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
    }
  },
  watch: {
    $route (to, from) {
      this.reviews = []
      this.$http.get('/reviews/' + this.type).then((response) => {
        this.reviews = response.data.data
      })
    }
  }
}
</script>

<!-- Add "scoped" attribute to limit CSS to this component only -->
<style scoped>
</style>
