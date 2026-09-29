import { useEffect } from 'react';

export function useDocumentTitle(page) {
  useEffect(() => {
    document.title = `VSRT-Lang — ${page}`;
  }, [page]);
}
