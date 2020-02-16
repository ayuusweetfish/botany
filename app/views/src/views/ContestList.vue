<template>
  <div>
    <v-container>
      <v-row>
        <div class="ml-4">当前共有{{total}}场比赛</div>
      </v-row>
      <v-row>
        <v-col
          :cols="12" :md="6"
          v-for="(item, index) in contests" :key="index"
        >
          <v-card>
            <v-card-title>{{item.title}}</v-card-title>
            <v-card-subtitle>{{item.time}}</v-card-subtitle>
            <v-card-text>
              <div>{{item.desc}}</div>
            </v-card-text>
          </v-card>
        </v-col>
      </v-row>
    </v-container>
  </div>
</template>

<script>
export default {
  name: 'ContestList',
  data: () => ({
    total: 0,
    contests: []
  }),
  mounted () {
    this.getContestList()
  },
  methods: {
    getContestList () {
      this.$axios.get('/contest/list').then(res => {
        this.contests = res.data
        this.total = res.data.length
        this.contests.forEach(item => {
          const start = this.$functions.dateTimeString(item.start_time)
          const end = this.$functions.dateTimeString(item.end_time)
          item.time = start + ' TO ' + end
        })
      })
    }
  }
}
</script>

<style>

</style>
