<template>
  <v-container fluid>
    <v-row align="center">
      <v-col :cols="12" :md="6">
        <v-img src="../assets/logo-l.png" :height="$vuetify.breakpoint.mdAndUp? '640px' : '280px'" contain></v-img>
      </v-col>
      <v-col :cols="12" :md="6">
        <router-view/>
      </v-col>
    </v-row>
  </v-container>
</template>

<script>
export default {
  name: 'Register'
}
</script>

<style>

</style>
