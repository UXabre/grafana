import { urlUtil } from './url';

describe('toUrlParams', () => {
  it('should encode object properties as url parameters', () => {
    const url = urlUtil.toUrlParams({
      server: 'backend-01',
      hasSpace: 'has space',
      many: ['1', '2', '3'],
      true: true,
      number: 20,
      isNull: null,
      isUndefined: undefined,
    });
    expect(url).toBe('server=backend-01&hasSpace=has%20space&many=1&many=2&many=3&true&number=20&isNull=&isUndefined=');
  });
});

describe('toUrlParams', () => {
  it('should encode the same way as angularjs', () => {
    const url = urlUtil.toUrlParams({
      server: ':@',
    });
    expect(url).toBe('server=:@');
  });
});

describe('parseKeyValue', () => {
  it('should parse url search params to object', () => {
    const obj = urlUtil.parseKeyValue('param=value&param2=value2&kiosk');
    expect(obj).toEqual({ param: 'value', param2: 'value2', kiosk: true });
  });

  it('should parse same url key multiple times to array', () => {
    const obj = urlUtil.parseKeyValue('servers=A&servers=B');
    expect(obj).toEqual({ servers: ['A', 'B'] });
  });
});
