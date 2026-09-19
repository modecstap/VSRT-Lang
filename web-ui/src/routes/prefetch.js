export function prefetchAccount() {
  import(/* webpackPrefetch: true */ '../pages/AccountPage/AccountPage');
  import(/* webpackPrefetch: true */ '../pages/AccountPage/components/ProfileTab/ProfileTab');
}

export function prefetchSessionTab() {
  import(/* webpackPrefetch: true */ '../pages/AccountPage/components/SessionTab/SessionTab');
}
