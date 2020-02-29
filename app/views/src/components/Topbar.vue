<template>
  <div>
    <div class="d-flex align-center" v-if="$vuetify.breakpoint.lgAndUp">
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
        <v-tab
          :to="`/contest/${$route.params.cid}/script`"
          v-if="checkAuth('moderator')"
        >脚本操作</v-tab>
        <v-tab
          :to="`/contest/${$route.params.cid}/judge#submit`"
          v-if="checkAuth('moderator')"
        >设置裁判</v-tab>
        <v-tab
          :to="`/contest/${$route.params.cid}/edit`"
          v-if="checkAuth('moderator')"
        >赛事编辑</v-tab>
        <v-tab
          :to="`/contest/${$route.params.cid}/participant#submit`"
          v-if="checkAuth('imIn')"
        >我的提交</v-tab>
        <v-tab :to="`/contest/${$route.params.cid}/ranklist`">选手排行</v-tab>
        <v-tab :to="`/contest/${$route.params.cid}/match`">对局列表</v-tab>
        <v-tab :to="`/contest/${$route.params.cid}/submission`">提交列表</v-tab>
      </v-tabs>
    </div>
    <div class="d-flex align-center" v-else-if="$vuetify.breakpoint.md">
      <div>
        <v-chip
          class="headline flex-nowrap"
          label
          color="white"
          large
          :to="`/contest/${$route.params.cid}/main`"
        >{{$store.state.cname}}</v-chip>
      </div>
      <div></div>
      <div>
        <v-menu
          v-model="menu"
          offset-y
          transition="slide-y-transition"
        >
          <template v-slot:activator="{ on }">
            <v-btn v-on="on" color="primary" icon><v-icon>mdi-menu</v-icon></v-btn>
          </template>
          <v-list light>
            <v-list-item
              :to="`/contest/${$route.params.cid}/script`"
              v-if="checkAuth('moderator')"
              active-class="primary--text"
            >
              <v-list-item-icon><v-icon>mdi-script-outline</v-icon></v-list-item-icon>
              <v-list-item-content>
                <v-list-item-title>脚本操作</v-list-item-title>
              </v-list-item-content>
            </v-list-item>
            <v-list-item
              :to="`/contest/${$route.params.cid}/judge#submit`"
              v-if="checkAuth('moderator')"
              active-class="primary--text"
            >
              <v-list-item-icon><v-icon>mdi-account-cog-outline</v-icon></v-list-item-icon>
              <v-list-item-content>
                <v-list-item-title>设置裁判</v-list-item-title>
              </v-list-item-content>
            </v-list-item>
            <v-list-item
              :to="`/contest/${$route.params.cid}/edit`"
              v-if="checkAuth('moderator')"
              active-class="primary--text"
            >
              <v-list-item-icon><v-icon>mdi-settings-outline</v-icon></v-list-item-icon>
              <v-list-item-content>
                <v-list-item-title>赛事编辑</v-list-item-title>
              </v-list-item-content>
            </v-list-item>
            <v-list-item
              :to="`/contest/${$route.params.cid}/participant#submit`"
              v-if="checkAuth('imIn')"
            >
              <v-list-item-icon><v-icon>mdi-code-tags</v-icon></v-list-item-icon>
              <v-list-item-content>
                <v-list-item-title>我的提交</v-list-item-title>
              </v-list-item-content>
            </v-list-item>
            <v-list-item
              :to="`/contest/${$route.params.cid}/ranklist`"
              active-class="primary--text"
            >
              <v-list-item-icon><v-icon>mdi-trophy</v-icon></v-list-item-icon>
              <v-list-item-content>
                <v-list-item-title>选手排行</v-list-item-title>
              </v-list-item-content>
            </v-list-item>
            <v-list-item
              :to="`/contest/${$route.params.cid}/match`"
              active-class="primary--text"
            >
              <v-list-item-icon><v-icon>mdi-format-list-checkbox</v-icon></v-list-item-icon>
              <v-list-item-content>
                <v-list-item-title>对局列表</v-list-item-title>
              </v-list-item-content>
            </v-list-item>
            <v-list-item
              :to="`/contest/${$route.params.cid}/submission`"
              active-class="primary--text"
            >
              <v-list-item-icon><v-icon>mdi-json</v-icon></v-list-item-icon>
              <v-list-item-content>
                <v-list-item-title>提交列表</v-list-item-title>
              </v-list-item-content>
            </v-list-item>
          </v-list>
        </v-menu>
      </div>
    </div>
    <div v-else>
      <v-menu
        v-model="menu"
        offset-y
        transition="slide-y-transition"
      >
        <template v-slot:activator="{ on }">
          <v-btn v-on="on" color="primary" icon><v-icon>mdi-menu</v-icon></v-btn>
        </template>
        <v-list light>
          <v-list-item
            :to="`/contest/${$route.params.cid}/main`"
            active-class="primary--text"
          >
            <v-list-item-content>
              <v-list-item-title class="title">{{$store.state.cname}}</v-list-item-title>
            </v-list-item-content>
          </v-list-item>
          <v-list-item
            :to="`/contest/${$route.params.cid}/script`"
            v-if="checkAuth('moderator')"
            active-class="primary--text"
          >
            <v-list-item-icon><v-icon>mdi-script-outline</v-icon></v-list-item-icon>
            <v-list-item-content>
              <v-list-item-title>脚本操作</v-list-item-title>
            </v-list-item-content>
          </v-list-item>
          <v-list-item
            :to="`/contest/${$route.params.cid}/judge#submit`"
            v-if="checkAuth('moderator')"
            active-class="primary--text"
          >
            <v-list-item-icon><v-icon>mdi-account-cog-outline</v-icon></v-list-item-icon>
            <v-list-item-content>
              <v-list-item-title>设置裁判</v-list-item-title>
            </v-list-item-content>
          </v-list-item>
          <v-list-item
            :to="`/contest/${$route.params.cid}/edit`"
            v-if="checkAuth('moderator')"
            active-class="primary--text"
          >
            <v-list-item-icon><v-icon>mdi-settings-outline</v-icon></v-list-item-icon>
            <v-list-item-content>
              <v-list-item-title>赛事编辑</v-list-item-title>
            </v-list-item-content>
          </v-list-item>
          <v-list-item
            :to="`/contest/${$route.params.cid}/participant#submit`"
            v-if="checkAuth('imIn')"
          >
            <v-list-item-icon><v-icon>mdi-code-tags</v-icon></v-list-item-icon>
            <v-list-item-content>
              <v-list-item-title>我的提交</v-list-item-title>
            </v-list-item-content>
          </v-list-item>
          <v-list-item
            :to="`/contest/${$route.params.cid}/ranklist`"
            active-class="primary--text"
          >
            <v-list-item-icon><v-icon>mdi-trophy</v-icon></v-list-item-icon>
            <v-list-item-content>
              <v-list-item-title>选手排行</v-list-item-title>
            </v-list-item-content>
          </v-list-item>
          <v-list-item
            :to="`/contest/${$route.params.cid}/match`"
            active-class="primary--text"
          >
            <v-list-item-icon><v-icon>mdi-format-list-checkbox</v-icon></v-list-item-icon>
            <v-list-item-content>
              <v-list-item-title>对局列表</v-list-item-title>
            </v-list-item-content>
          </v-list-item>
          <v-list-item
            :to="`/contest/${$route.params.cid}/submission`"
            active-class="primary--text"
          >
            <v-list-item-icon><v-icon>mdi-json</v-icon></v-list-item-icon>
            <v-list-item-content>
              <v-list-item-title>提交列表</v-list-item-title>
            </v-list-item-content>
          </v-list-item>
        </v-list>
      </v-menu>
    </div>
  </div>
</template>

<script>
export default {
  data: () => ({
    settingMenu: false,
    tab: null,
    menu: false
  }),
  methods: {
    checkAuth (type) {
      return this.$store.state.myrole === this.$consts.role[type]
    }
  }
}
</script>

<style>

</style>
