resource "aws_s3_bucket" "package_store" {
  bucket = var.package_store_s3_bucket_name
  tags   = var.tags
}

resource "aws_s3_bucket_server_side_encryption_configuration" "package_store" {
  bucket = aws_s3_bucket.package_store.id
  rule {
    apply_server_side_encryption_by_default {
      sse_algorithm = "AES256"
    }
  }
}

resource "aws_s3_bucket_public_access_block" "package_store" {
  bucket                  = aws_s3_bucket.package_store.bucket
  block_public_acls       = true
  block_public_policy     = true
  ignore_public_acls      = true
  restrict_public_buckets = true
}

resource "aws_s3_bucket_ownership_controls" "package_store" {
  bucket = aws_s3_bucket.package_store.id
  rule {
    object_ownership = "BucketOwnerEnforced"
  }
}
