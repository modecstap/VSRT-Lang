import { renderHook } from '@testing-library/react';
import { useDocumentTitle } from './useDocumentTitle';

test('sets tab title and follows page changes', () => {
  const { rerender } = renderHook(({ page }) => useDocumentTitle(page), { initialProps: { page: 'Login' } });
  expect(document.title).toBe('VSRT-Lang — Login');
  rerender({ page: 'Register' });
  expect(document.title).toBe('VSRT-Lang — Register');
});
