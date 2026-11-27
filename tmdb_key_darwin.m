//go:build darwin && cgo
#import <Foundation/Foundation.h>
#import <Security/Security.h>
static NSMutableDictionary *tmdbQuery(const char *account){
 return [NSMutableDictionary dictionaryWithDictionary:@{
  (id)kSecClass:(id)kSecClassGenericPassword,
  (id)kSecAttrService:@"es.cristianbarallobre.selectorcine365.tmdb",
  (id)kSecAttrAccount:[NSString stringWithUTF8String:account]
 }];
}
int cineTMDBKeyRead(const char *account,void **out,int *length){
 @autoreleasepool {
  NSMutableDictionary *q=tmdbQuery(account);
  q[(id)kSecReturnData]=@YES;
  CFTypeRef item=NULL;
  OSStatus status=SecItemCopyMatching((CFDictionaryRef)q,&item);
  if(status!=errSecSuccess)return status;
  NSData *d=(NSData *)item;
  if(!d.length||d.length>32768){CFRelease(item);return errSecDecode;}
  *length=(int)d.length;*out=malloc(d.length);
  if(*out)memcpy(*out,d.bytes,d.length);
  CFRelease(item);
  return *out?errSecSuccess:errSecAllocate;
 }
}
int cineTMDBKeySave(const char *account,const void *data,int length){
 @autoreleasepool {
  if(length<=0||length>32768)return errSecParam;
  NSMutableDictionary *q=tmdbQuery(account);
  NSDictionary *attrs=@{(id)kSecValueData:[NSData dataWithBytes:data length:length]};
  OSStatus status=SecItemUpdate((CFDictionaryRef)q,(CFDictionaryRef)attrs);
  if(status!=errSecItemNotFound)return status;
  [q addEntriesFromDictionary:attrs];
  return SecItemAdd((CFDictionaryRef)q,NULL);
 }
}
