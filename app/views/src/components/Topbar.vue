<template>
  <div class="d-flex align-center">
    <div>
      <v-chip
        class="headline flex-nowrap"
        label
        color="white"
        large
        :to="`/contest/${$route.params.cid}/main`"
      >{{$store.state.cname}}</v-chip>
    </div>
    <v-tabs v-model="tab" height="64" grow>
      <v-tab :to="`/contest/${$route.params.cid}/participant#submit`">我的提交</v-tab>
      <v-tab :to="`/contest/${$route.params.cid}/ranklist`">选手排行</v-tab>
      <v-tab :to="`/contest/${$route.params.cid}/match`">对局列表</v-tab>
      <v-tab :to="`/contest/${$route.params.cid}/submission`">提交列表</v-tab>
      <v-tab :to="`/contest/${$route.params.cid}/edit`">赛事编辑</v-tab>
    </v-tabs>
  </div>
  <!-- <div class="d-flex align-center">
    <v-btn
      text
      class="subtitle-1 font-weight-bold"
      color="secondary"
      :to="`/contest/${$route.params.cid}/main`"
    >{{$store.state.cname}}</v-btn>
    <v-menu
      v-if="$store.state.myrole===$consts.role.moderator"
      v-model="settingMenu"
      offset-y
    >
      <template v-slot:activator="{ on }">
        <v-btn
          text
          v-on="on"
          color="secondary"
          class="subtitle-1"
        >赛事设置</v-btn>
      </template>
      <v-list>
        <v-list-item>
          <v-list-item-title
            class="subtitle-1"
          >
            <v-icon>mdi-settings-outline</v-icon>
            编辑信息
          </v-list-item-title>
        </v-list-item>
      </v-list>
    </v-menu>
    <v-btn text color="secondary" class="subtitle-1">确实</v-btn>
  </div> -->
</template>

<script>
export default {
  data: () => ({
    settingMenu: false,
    tab: null
  })
}
</script>

<style>

</style>
