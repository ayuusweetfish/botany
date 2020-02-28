<template>
  <v-container fluid>
    <v-row justify="center">
      <v-col :cols="12" :md="10" :lg="8">
        <v-card outlined style="border: none">
          <v-card-title class="title">新建赛事</v-card-title>
          <v-card-text>
            <div class="body-1">在账号{{$store.state.handle}}下添加一场赛事</div>
          </v-card-text>
        </v-card>
        <v-divider class="mb-2"></v-divider>
        <contest-editor
          :cid="$route.params.cid"
          mode="create"
        ></contest-editor>
      </v-col>
    </v-row>
  </v-container>
</template>

<script>
import ContestEditor from '../components/ContestEditor.vue'
export default {
  name: 'ContestEdit',
  components: {
    'contest-editor': ContestEditor
  }
}
</script>
