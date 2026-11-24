//go:build darwin && cgo

#import <Foundation/Foundation.h>
#import <Security/Security.h>

static NSMutableDictionary *cineKeyQuery(const char *account) {
    return [NSMutableDictionary dictionaryWithDictionary:@{
        (id)kSecClass: (id)kSecClassGenericPassword,
        (id)kSecAttrService: @"es.cristianbarallobre.selectorcine365.onedrive",
        (id)kSecAttrAccount: [NSString stringWithUTF8String:account]
    }];
}
int cineKeychainRead(const char *account, void **out, int *length) {
    @autoreleasepool {
        NSMutableDictionary *q = cineKeyQuery(account);
        q[(id)kSecReturnData] = @YES;
        CFTypeRef item = NULL;
        OSStatus status = SecItemCopyMatching((CFDictionaryRef)q, &item);
        if (status != errSecSuccess) return status;
        NSData *d = (NSData *)item;
        if (!d.length || d.length > 32768) { CFRelease(item); return errSecDecode; }
        *length = (int)d.length;
        *out = malloc(d.length);
        if (*out) memcpy(*out, d.bytes, d.length);
        CFRelease(item);
        return *out ? errSecSuccess : errSecAllocate;
    }
}
int cineKeychainSave(const char *account, const void *data, int length) {
    @autoreleasepool {
        NSMutableDictionary *q = cineKeyQuery(account);
        NSDictionary *attributes = @{(id)kSecValueData:[NSData dataWithBytes:data length:length]};
        OSStatus status = SecItemUpdate((CFDictionaryRef)q, (CFDictionaryRef)attributes);
        if (status != errSecItemNotFound) return status;
        [q addEntriesFromDictionary:attributes];
        return SecItemAdd((CFDictionaryRef)q, NULL);
    }
}
int cineKeychainDelete(const char *account) {
    @autoreleasepool { return SecItemDelete((CFDictionaryRef)cineKeyQuery(account)); }
}
