<template>
  <v-container>
    <v-row justify="space-around">
      <v-btn v-for="n in names" :key="n" text color="primary"
        @click.prevent="loadReport(n)">{{ n }}
      </v-btn>
    </v-row>

    <v-row>
      <v-card>
        <v-data-table v-if="entries.length" :headers="headers" :items="entries"
          :items-per-page="30" class="elevation-1" :loading="loading"
          :sort-by="'lastSeen'" :sort-desc="true" calculate-widths dense>

          <template v-slot:[`item.lastSeen`]="{ item }">
            {{ item.lastSeen | prettyDateTime }}
          </template>

          <template v-slot:[`item.title`]="{ item }">
            <a class="text-decoration-none title-link " :href="item.url"
              target="_blank">{{
              item.title
              }}</a>
          </template>

        </v-data-table>
      </v-card>
    </v-row>

  </v-container>
</template>

<script>
import moment from 'moment'

export default {
  name: 'HelloWorld',
  props: ['password'],
  data: () => ({
    names: ['gar'],
    entries: [],
    loading: false,
    headers: [
      { text: 'Last Seen', value: 'lastSeen', width: '20%' },
      { text: 'Count', value: 'count', width: '10%' },
      { text: 'Title', value: 'title', width: '70%' },
    ]
  }),
  async created() {
    const response = await fetch(`https://instance-1.gmendiola.com/names`, {
      headers: {
        'X-Api-Key': this.password
      }
    })
    const data = await response.json()
    this.names = data?.names

    this.loadReport('gar')
  },
  methods: {
    async loadReport(name) {
      this.loading = true
      const response = await fetch(`https://instance-1.gmendiola.com/logs?name=${name}`, {
        headers: {
          'X-Api-Key': this.password
        }
      })
      const report = await response.json()
      this.entries = report?.entries
      this.loading = false
    },
  },
  filters: {
    prettyDateTime: function (date) {
      return moment(date).format('MM/DD/YYYY, h:mm:ss a');
    }
  }
}
</script>

<style>
.title-link {
  color: rgba(0, 0, 0, 0.87) !important;
  word-wrap: break-word;
  word-break: break-all;
  width: auto;
}
</style>